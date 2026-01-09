package service

import (
	"github.com/konsultin/project-goes-here/dto"
	"github.com/konsultin/project-goes-here/internal/svc-core/model"
	"github.com/konsultin/project-goes-here/internal/svc-core/pkg/svck"
	"github.com/konsultin/errk"
	logkOption "github.com/konsultin/logk/option"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

// getUserByXid retrieves user by XID
func (s *Service) getUserByXid(xid string) (*model.User, error) {
	user, err := s.repo.FindUserByXid(xid)
	if err != nil {
		s.log.Error("Failed to FindUserByXid", logkOption.Error(err))
		return nil, errk.Trace(err)
	}
	return user, nil
}

// getUserById retrieves user by ID
func (s *Service) getUserById(id int64) (*model.User, error) {
	user, err := s.repo.FindUserById(id)
	if err != nil {
		s.log.Error("Failed to FindUserById", logkOption.Error(err))
		return nil, errk.Trace(err)
	}
	return user, nil
}

// generateXid generates a new XID for user
func (s *Service) generateXid() string {
	return gonanoid.MustGenerate(svck.AlphaNumUpperCharSet, 12)
}

// composeControlStatusResult creates a Status DTO from ControlStatus_Enum
func composeControlStatusResult(statusId dto.ControlStatus_Enum) *dto.Status {
	name := "UNKNOWN"
	switch statusId {
	case dto.ControlStatus_ACTIVE:
		name = "ACTIVE"
	case dto.ControlStatus_INACTIVE:
		name = "INACTIVE"
	case dto.ControlStatus_LOCKED:
		name = "LOCKED"
	case dto.ControlStatus_PENDING:
		name = "PENDING"
	case dto.ControlStatus_DRAFT:
		name = "DRAFT"
	}
	return &dto.Status{
		Id:   int32(statusId),
		Name: name,
	}
}
