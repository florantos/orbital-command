package domain

import (
	"errors"
	"fmt"
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
	name = strings.TrimSpace(name)
	description = strings.TrimSpace(description)

	if name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if utf8.RuneCountInString(name) > 100 {
		return nil, fmt.Errorf("name cannot exceed 100 characters")
	}

	if description == "" {
		return nil, fmt.Errorf("description is required")
	}

	return &Module{
		Name:        name,
		Description: description,
		HealthState: HealthStateOperational,
	}, nil
}
