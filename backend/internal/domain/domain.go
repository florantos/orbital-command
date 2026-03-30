package domain

import (
	"sort"
	"strings"
)

type ValidationError struct {
	Fields map[string]string
}

func (e *ValidationError) Error() string {
	keys := make([]string, 0, len(e.Fields))
	for field := range e.Fields {
		keys = append(keys, field)
	}
	sort.Strings(keys)

	msgs := make([]string, 0, len(e.Fields))
	for _, key := range keys {
		msgs = append(msgs, key+": "+e.Fields[key])
	}
	return strings.Join(msgs, ", ")
}

type Capability string

const (
	CapabilityDocking             Capability = "docking"
	CapabilityNavigation          Capability = "navigation"
	CapabilityHullMonitoring      Capability = "hull-monitoring"
	CapabilityWaterRecycling      Capability = "water-recycling"
	CapabilityPowerGeneration     Capability = "power-generation"
	CapabilityPowerDistribution   Capability = "power-distubution"
	CapabilityThermalRegulation   Capability = "thermal-regulation"
	CapabilityAtmosphereRecycling Capability = "atmosphere-recycling"
)

var validCapabilities = map[Capability]struct{}{
	CapabilityDocking:             {},
	CapabilityNavigation:          {},
	CapabilityHullMonitoring:      {},
	CapabilityWaterRecycling:      {},
	CapabilityPowerGeneration:     {},
	CapabilityPowerDistribution:   {},
	CapabilityThermalRegulation:   {},
	CapabilityAtmosphereRecycling: {},
}
