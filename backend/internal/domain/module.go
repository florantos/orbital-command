package domain

import (
	"errors"
	"strings"
	"time"
	"unicode/utf8"
)

type HealthState string

const (
	HealthStateOperational  HealthState = "operational"
	HealthStateDegraded     HealthState = "degraded"
	HealthStateCritical     HealthState = "critical"
	HealthStateUnresponsive HealthState = "unresponsive"
	HealthStateOffline      HealthState = "offline"
)

type Module struct {
	ID          string
	Name        string
	Description string
	HealthState HealthState
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

var ErrDuplicateModuleName = errors.New("duplicate module name")

func NewModule(name, description string) (*Module, error) {
	m := &Module{
		Name:        name,
		Description: description,
		HealthState: HealthStateOperational,
	}

	err := m.Validate()
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (m *Module) Validate() error {
	ve := &ValidationError{
		Fields: make(map[string]string),
	}

	m.Name = strings.TrimSpace(m.Name)
	m.Description = strings.TrimSpace(m.Description)

	if m.Name == "" {
		ve.Fields["name"] = "name is required"
	} else if utf8.RuneCountInString(m.Name) > 100 {
		ve.Fields["name"] = "name cannot exceed 100 characters"
	}

	if m.Description == "" {
		ve.Fields["description"] = "description is required"
	}

	if len(ve.Fields) > 0 {
		return ve
	}
	return nil
}
