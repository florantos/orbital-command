package domain

import "strings"

type ValidationError struct {
	Fields map[string]string
}

func (e *ValidationError) Error() string {
	msgs := make([]string, 0, len(e.Fields))
	for field, msg := range e.Fields {
		msgs = append(msgs, field+": "+msg)
	}
	return strings.Join(msgs, ", ")
}
