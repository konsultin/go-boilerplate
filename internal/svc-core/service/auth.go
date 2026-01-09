package service

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/konsultin/project-goes-here/dto"
	"github.com/konsultin/project-goes-here/internal/svc-core/model"
	"github.com/konsultin/project-goes-here/internal/svc-core/pkg/httpk"
	"github.com/konsultin/project-goes-here/internal/svc-core/pkg/oauth/google"
	"github.com/konsultin/errk"
	logkOption "github.com/konsultin/logk/option"
	"github.com/konsultin/timek"
	"golang.org/x/crypto/bcrypt"
)

// LoginWithPassword authenticates user with identifier (email/phone/username) and password
// Requires anonymous session bearer token for authentication
func (s *Service) LoginWithPassword(payload *dto.LoginPassword_Payload) (*dto.CreateUserSession_Result_Data, error) {
	// Verify anonymous session token first
	if err := s.verifyAnonymousSession(); err != nil {
		return nil, err
	}

	// Normalize identifier (lowercase for email)
	identifier := strings.TrimSpace(payload.Identifier)
	if strings.Contains(identifier, "@") {
		identifier = strings.ToLower(identifier)
	}

	// Find user by identifier
	user, err := s.repo.FindUserByIdentifier(identifier)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.log.Warnf("User not found for identifier: %s", identifier)
			return nil, httpk.UnauthorizedError
		}
		s.log.Error("Failed to find user", logkOption.Error(err))
		return nil, errk.Trace(err)
	}

	// Check user status
	if user.StatusId != dto.ControlStatus_ACTIVE {
		s.log.Warnf("User account is not active. UserId=%d Status=%d", user.Id, user.StatusId)
		return nil, httpk.ForbiddenError
	}

	// Find password credential for this user
	credential, err := s.repo.FindCredentialByKey(dto.AuthProvider_PASSWORD, identifier)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.log.Warnf("No password credential found for identifier: %s", identifier)
			return nil, httpk.UnauthorizedError
		}
		s.log.Error("Failed to find credential", logkOption.Error(err))
		return nil, errk.Trace(err)
	}

	// Verify password
	if !credential.CredentialSecret.Valid {
		s.log.Warnf("Credential has no password set. CredentialId=%d", credential.Id)
		return nil, httpk.UnauthorizedError
	}

	err = bcrypt.CompareHashAndPassword([]byte(credential.CredentialSecret.String), []byte(payload.Password))
	if err != nil {
		s.log.Warnf("Invalid password for identifier: %s", identifier)
		return nil, httpk.UnauthorizedError
	}

	// Create user session
	return s.CreateUserSession(user, dto.AuthProvider_PASSWORD, payload.Device, time.Now())
}

// LoginWithGoogle authenticates user with Google OAuth
// Requires anonymous session bearer token for authentication
func (s *Service) LoginWithGoogle(payload *dto.LoginOAuth_Payload) (*dto.CreateUserSession_Result_Data, error) {
	// Verify anonymous session token first
	if err := s.verifyAnonymousSession(); err != nil {
		return nil, err
	}

	// Create Google provider
	provider := google.NewProvider(s.config.GoogleClientID)

	// Verify Google token
	userInfo, err := provider.VerifyToken(s.ctx, payload.IdToken)
	if err != nil {
		s.log.Error("Failed to verify Google token", logkOption.Error(err))
		return nil, httpk.UnauthorizedError.Wrap(err)
	}

	// Find existing credential for Google + user ID
	credential, err := s.repo.FindCredentialByKey(dto.AuthProvider_GOOGLE, userInfo.ProviderId)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		s.log.Error("Failed to find credential", logkOption.Error(err))
		return nil, errk.Trace(err)
	}

	var user *model.User

	if credential != nil {
		// User exists, get user record
		user, err = s.getUserById(credential.UserId)
		if err != nil {
			s.log.Error("Failed to find user", logkOption.Error(err))
			return nil, errk.Trace(err)
		}
	} else {
		// New OAuth user - create user and credential
		user, err = s.createGoogleUser(userInfo)
		if err != nil {
			s.log.Error("Failed to create Google user", logkOption.Error(err))
			return nil, errk.Trace(err)
		}
	}

	// Check user status
	if user.StatusId != dto.ControlStatus_ACTIVE {
		s.log.Warnf("User account is not active. UserId=%d Status=%d", user.Id, user.StatusId)
		return nil, httpk.ForbiddenError
	}

	// Create user session
	return s.CreateUserSession(user, dto.AuthProvider_GOOGLE, payload.Device, time.Now())
}

// verifyAnonymousSession checks if request has valid anonymous session bearer token
func (s *Service) verifyAnonymousSession() error {
	// Get bearer token from context
	token, ok := s.ctx.Value(httpk.BearerToken).(string)
	if !ok || token == "" {
		s.log.Warn("Missing bearer token for login request")
		return httpk.UnauthorizedError
	}

	// Verify the JWT token (no audience check for anonymous session validation)
	jwtAdapter := s.NewJwtAdapter()
	claims, err := jwtAdapter.ValidateWithoutAudience(token)
	if err != nil {
		s.log.Error("Failed to validate anonymous session token", logkOption.Error(err))
		return httpk.UnauthorizedError.Wrap(err)
	}

	// Check if it's an anonymous session (ANONYMOUS_USER or ANONYMOUS_ADMIN)
	subjectType := dto.Role_Enum(claims.Ent)
	if subjectType != dto.Role_ANONYMOUS_USER && subjectType != dto.Role_ANONYMOUS_ADMIN {
		s.log.Warnf("Invalid session type for login. Expected anonymous, got: %d", subjectType)
		return httpk.UnauthorizedError
	}

	return nil
}

// createGoogleUser creates a new user from Google user info
func (s *Service) createGoogleUser(userInfo *google.UserInfo) (*model.User, error) {
	now := timek.Now()

	// Create user
	user := &model.User{
		BaseField: model.NewBaseFieldFromModel(s.subject),
		Xid:       s.generateXid(),
		FullName:  userInfo.Name,
		Email:     sql.NullString{String: userInfo.Email, Valid: userInfo.Email != ""},
		Avatar:    sql.NullString{String: userInfo.Picture, Valid: userInfo.Picture != ""},
		StatusId:  dto.ControlStatus_ACTIVE,
	}

	// Insert user
	err := s.repo.InsertUser(user)
	if err != nil {
		return nil, errk.Trace(err)
	}

	// Create credential
	credential := &model.UserCredential{
		UserId:         user.Id,
		AuthProviderId: dto.AuthProvider_GOOGLE,
		CredentialKey:  userInfo.ProviderId,
		IsVerified:     userInfo.EmailVerified,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if userInfo.EmailVerified {
		credential.VerifiedAt = sql.NullTime{Time: time.Now(), Valid: true}
	}

	err = s.repo.InsertUserCredential(credential)
	if err != nil {
		return nil, errk.Trace(err)
	}

	return user, nil
}
