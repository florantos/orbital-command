package domain

import (
	"errors"
	"strings"
	"time"
)

type Role string

const (
	RoleEngineer   Role = "engineer"
	RoleTechnician Role = "technician"
	RoleSpecialist Role = "specialist"
)

var validRoles = map[Role]struct{}{
	RoleEngineer:   {},
	RoleTechnician: {},
	RoleSpecialist: {},
}

type CrewMember struct {
	ID             string
	Name           string
	Role           Role
	Qualifications []Capability
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

var ErrDuplicateCrewMemberName = errors.New("duplicate crew member name")

func NewCrewMember(name string, role Role, qualifications []Capability) (*CrewMember, error) {
	cm := &CrewMember{
		Name:           strings.TrimSpace(name),
		Role:           role,
		Qualifications: qualifications,
	}
	err := cm.Validate()
	if err != nil {
		return nil, err
	}
	return cm, nil
}

func (cm *CrewMember) Validate() error {
	ve := &ValidationError{
		Fields: make(map[string]string),
	}

	if cm.Name == "" {
		ve.Fields["name"] = "name is required"
	}

	if _, ok := validRoles[cm.Role]; !ok {
		ve.Fields["role"] = string(cm.Role) + " is not a valid role"
	}

	if len(cm.Qualifications) == 0 {
		ve.Fields["qualifications"] = "at least one qualification is required"
	} else {
		for _, q := range cm.Qualifications {
			if _, ok := validCapabilities[q]; !ok {
				ve.Fields["qualifications"] = string(q) + " is not a valid qualification"
				break
			}
		}
	}

	if len(ve.Fields) > 0 {
		return ve
	}
	return nil
}
