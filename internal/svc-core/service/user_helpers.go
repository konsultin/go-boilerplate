package service

import (
	"github.com/konsultin/project-goes-here/dto"
	"github.com/konsultin/project-goes-here/internal/svc-core/model"
	"github.com/konsultin/project-goes-here/libs/errk"
	logkOption "github.com/konsultin/project-goes-here/libs/logk/option"
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
