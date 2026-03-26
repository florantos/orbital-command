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
