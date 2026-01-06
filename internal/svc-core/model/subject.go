package model

import "github.com/konsultin/project-goes-here/dto"

type Subject struct {
	Id       string `json:"id"`
	Role     string `json:"role"`
	FullName string `json:"fullName"`
}

func NewSubject(d *dto.Subject) *Subject {
	if d == nil {
		return &Subject{}
	}

	return &Subject{
		Id:       d.Id,
		Role:     d.Role,
		FullName: d.FullName,
	}
}

func ToSubjectResult(s *Subject) *dto.Subject {
	if s == nil {
		return nil
	}

	return &dto.Subject{
		Id:       s.Id,
		Role:     s.Role,
		FullName: s.FullName,
	}
}
