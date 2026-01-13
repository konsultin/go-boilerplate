package service

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/go-konsultin/errk"
	logkOption "github.com/go-konsultin/logk/option"
	"github.com/konsultin/project-goes-here/dto"
	specErr "github.com/konsultin/project-goes-here/internal/errors"
	"github.com/konsultin/project-goes-here/internal/svc-core/constant"
	"github.com/konsultin/project-goes-here/internal/svc-core/model"
	"github.com/konsultin/project-goes-here/internal/svc-core/pkg/httpk"
	unaryHttpk "github.com/konsultin/project-goes-here/internal/svc-core/pkg/httpk/unary"
	"github.com/konsultin/project-goes-here/internal/svc-core/pkg/svck"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"golang.org/x/crypto/bcrypt"
)

func (s *Service) CreateAnonymousUserSession(payload *unaryHttpk.BasicAuth, clientTypeId dto.Role_Enum) (*dto.CreateUserSession_Result, error) {
	clientAuth, err := s.repo.FindClientAuthByClientId(payload.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.log.Errorf("clientAuth is not found. Username = %s", payload.Username)
			return nil, specErr.InvalidCredentials
		}
		s.log.Error("Failed to FindClientAuthByClientId", logkOption.Error(err))
		return nil, errk.Trace(err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(clientAuth.Options.ClientSecret), []byte(payload.Password))
	if err != nil {
		s.log.Error("Failed to compare the password", logkOption.Error(err))
		return nil, specErr.InvalidCredentials.Wrap(err).Trace()
	}

	if clientAuth.ClientTypeId != clientTypeId {
		s.log.Error("Invalid clientTypeId")
		return nil, specErr.InvalidClientType.Trace()
	}

	var subjectType int32
	switch clientAuth.ClientTypeId {
	case dto.Role_ANONYMOUS_ADMIN, dto.Role_ANONYMOUS_USER:
		subjectType = int32(clientAuth.ClientTypeId)
	default:
		return nil, fmt.Errorf("invalid Client Type. clientTypeId = %v", clientAuth.ClientTypeId)
	}

	rolePrivileges, fErr := s.repo.FindRolePrivilegeByRoleId(subjectType)
	if fErr != nil {
		s.log.Error("Failed to FindRolePrivilegeByRoleId", logkOption.Error(fErr))
		return nil, errk.Trace(fErr)
	}

	var audience []string
	for _, val := range rolePrivileges {
		audience = append(audience, val.Privilege.Xid)
	}

	jwtAdapter := s.NewJwtAdapter()
	session, err := jwtAdapter.Issue(IssueJwtPayload{
		Subject:     clientAuth.ClientId,
		Audience:    audience,
		Lifetime:    clientAuth.Options.TokenLifetime,
		SessionId:   gonanoid.MustGenerate(svck.AlphaNumUpperCharSet, 6),
		SubjectType: subjectType,
	})

	if err != nil {
		s.log.Error("Failed to issue jwt payload", logkOption.Error(err))
		return nil, errk.Trace(err)
	}

	return &dto.CreateUserSession_Result{
		Session: session,
		Scopes:  audience,
	}, nil
}

func (s *Service) RefreshUserSession(payload *dto.RefreshSession_Payload) (*dto.CreateUserSession_Result_Data, error) {
	// get the token
	token, ok := s.ctx.Value(httpk.BearerToken).(string)
	if !ok {
		token = ""
	}

	// Verify jwt
	jwtAdapter := s.NewJwtAdapter()
	jwtToken, err := jwtAdapter.Validate(token, &dto.ValidateJwt_Payload{
		Audience: []string{constant.PrivilegeRefreshUserToken},
	})
	if err != nil {
		s.log.Error("Failed to validate jwt payload", logkOption.Error(err))
		return nil, httpk.UnauthorizedError.Wrap(err).Trace()
	}

	// get user by xid from token
	user, err := s.getUserByXid(jwtToken.Sub)
	if err != nil {
		return nil, httpk.UnauthorizedError.Wrap(err).Trace()
	}

	// get session by xid from token
	session, err := s.repo.FindSessionByXid(jwtToken.Jti)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.log.Error("No Session found", logkOption.Error(err))
			return nil, httpk.UnauthorizedError.Wrap(err).Trace()
		}
		s.log.Error("Failed to FindSessionByXid", logkOption.Error(err))
		return nil, errk.Trace(err)
	}
	// Check Auth Session status
	if err = s.isValidAuthSession(session); err != nil {
		return nil, errk.Trace(err)
	}
	// Create new user session
	data, err := s.CreateUserSession(user, session.AuthProviderId, payload.Device, time.Now())
	if err != nil {
		s.log.Error("Failed to CreateUserSession", logkOption.Error(err))
		return nil, errk.Trace(err)
	}

	// Delete previous session
	err = s.repo.DeleteSessionByXid(session.Xid)
	if err != nil {
		s.log.Error("Unable to delete previous session. SessionXid=%s", logkOption.Error(err), logkOption.Format(session.Xid))
		return nil, errk.Trace(err)
	}

	return data, nil
}

func (s *Service) isValidAuthSession(session *model.AuthSession) error {
	switch session.StatusId {
	case dto.ControlStatus_ACTIVE:
		// Do nothing
	case dto.ControlStatus_LOCKED:
		// If session status is marked as Locked, then return E_AUTH_2 to trigger Client to Refresh Session
		s.log.Warnf("Session is marked as expired to trigger client to refresh session. Id=%d SubjectId=%s SubjectTypeId=%d",
			session.Id, session.SubjectId, session.SubjectTypeId)
		return constant.CurrentAuthSessionExpired
	case dto.ControlStatus_INACTIVE:
		if s.config.FeatureFlagSingleDevice {
			// Delete current session
			err := s.DeleteSession(session.Xid)
			if err != nil {
				return errk.Trace(err)
			}
			// if session status is marked as Inactive, then return E_AUTH_3
			s.log.Warnf("Session is marked as inactive to trigger client to re-login. Id=%d SubjectId=%s", session.Id, session.SubjectId)
			return constant.LoginDetectedAnotherDevice
		}
	default:
		s.log.Warnf("unexpected auth session status. StatusId=%s", session.StatusId)
		return httpk.UnauthorizedError
	}
	return nil
}

func (s *Service) DeleteSession(xid string) error {
	err := s.repo.DeleteSessionByXid(xid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return constant.ResourceNotFound
		}
		s.log.Error("Failed to delete session", logkOption.Error(err))
		return errk.Trace(err)
	}
	return nil
}

func (s *Service) CreateUserSession(user *model.User, authProviderId dto.AuthProvider_Enum, device *dto.DeviceSession, t time.Time) (*dto.CreateUserSession_Result_Data, error) {
	// Get user privileges
	subjectType := int32(dto.Role_USER)
	rolePrivileges, err := s.repo.FindRolePrivilegeByRoleId(subjectType)
	if err != nil {
		s.log.Error("Failed to FindRolePrivilegeByRoleId on CreateUserSession", logkOption.Error(err))
		return nil, errk.Trace(err)
	}
	var audience []string
	for _, val := range rolePrivileges {
		audience = append(audience, val.Privilege.Xid)
	}

	// Get created At
	createdAt := sql.NullTime{Time: t, Valid: true}
	sessionId := gonanoid.MustGenerate(svck.AlphaNumUpperCharSet, 10)
	jwtAdapter := s.NewJwtAdapter()
	// Issue the JWT for Access Token
	accessSession, err := jwtAdapter.Issue(IssueJwtPayload{
		SessionId:   sessionId,
		Subject:     user.Xid,
		Audience:    audience,
		Lifetime:    s.config.UserSessionLifetime,
		SubjectType: subjectType,
		CreatedAt:   createdAt,
	})

	// Issue the JWT for Refresh Token
	refreshSession, err := jwtAdapter.Issue(IssueJwtPayload{
		SessionId:   sessionId,
		Subject:     user.Xid,
		Audience:    []string{constant.PrivilegeRefreshUserToken},
		Lifetime:    s.config.UserSessionRefreshLifetime,
		SubjectType: subjectType,
		CreatedAt:   createdAt,
	})

	// Init baseField
	baseField := model.NewBaseFieldFromModel(s.subject)

	// FCM Token
	notificationToken := sql.NullString{}
	if device.NotificationToken != "" {
		notificationToken.Valid = true
		notificationToken.String = device.NotificationToken
	}

	authSession := &model.AuthSession{
		BaseField:        baseField,
		Xid:              sessionId,
		SubjectId:        user.Xid,
		SubjectTypeId:    dto.Role_Enum(subjectType),
		AuthProviderId:   authProviderId,
		DevicePlatformId: device.DevicePlatformId,
		DeviceId:         device.DeviceId,
		Device: &model.AuthSessionDevice{
			DeviceId:         device.DeviceId,
			DevicePlatformId: device.DevicePlatformId,
			ClientIp:         "",
		},
		NotificationChannelId: device.NotificationChannelId,
		NotificationToken:     notificationToken,
		ExpiredAt:             time.Unix(accessSession.ExpiredAt, 0),
		StatusId:              dto.ControlStatus_ACTIVE,
	}

	// Persist
	err = s.repo.InsertAuthSession(authSession)
	if err != nil {
		s.log.Error("Failed to InsertAuthSession", logkOption.Error(err))
		return nil, errk.Trace(err)
	}

	return &dto.CreateUserSession_Result_Data{
		User:           s.mustComposeUserResult(user),
		AccessSession:  accessSession,
		RefreshSession: refreshSession,
		AccessScopes:   audience,
	}, nil
}

func (s *Service) mustComposeUserResult(m *model.User) *dto.User {
	return &dto.User{
		Id:         m.Id,
		Xid:        m.Xid,
		Phone:      m.Phone.String,
		FullName:   m.FullName,
		Email:      m.Email.String,
		Age:        m.Age.String,
		Avatar:     s.mustComposeFileResult(dto.AssetType_USER_AVATAR, m.Avatar.String),
		Status:     composeControlStatusResult(m.StatusId),
		ModifiedBy: model.ToSubjectResult(m.ModifiedBy),
		CreatedAt:  m.CreatedAt.ToTime().Unix(),
		UpdatedAt:  m.UpdatedAt.ToTime().Unix(),
		Version:    m.Version,
	}
}
