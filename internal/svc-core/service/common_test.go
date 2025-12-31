package service

import (
	"testing"

	"github.com/konsultin/project-goes-here/dto"
	"github.com/konsultin/project-goes-here/internal/svc-core/repository"
)

func TestService_WithSubject(t *testing.T) {
	// Setup
	repo := &repository.Repository{} // Mock repo
	svc := NewService(repo)

	// Define test subject
	dtoSubject := &dto.Subject{
		Id:       "user-123",
		Role:     "admin",
		FullName: "Test User",
	}

	// Test WithSubject
	svcWithSubject := svc.WithSubject(dtoSubject)

	// Verify new service instance is different (pointer check)
	if svc == svcWithSubject {
		t.Error("WithSubject should return a new pointer")
	}

	// Verify subject is set in new instance
	if svcWithSubject.subject == nil {
		t.Fatal("Subject should be set in new instance")
	}
	if svcWithSubject.subject.Id != dtoSubject.Id {
		t.Errorf("Expected subject ID %s, got %s", dtoSubject.Id, svcWithSubject.subject.Id)
	}
	if svcWithSubject.subject.Role != dtoSubject.Role {
		t.Errorf("Expected subject Role %s, got %s", dtoSubject.Role, svcWithSubject.subject.Role)
	}

	// Verify original instance is untouched
	if svc.subject != nil {
		t.Error("Original service instance should not have subject set")
	}

	// Verify dependencies are preserved (shared)
	if svcWithSubject.repo != svc.repo {
		t.Error("Repository reference should be preserved/shared")
	}
}
