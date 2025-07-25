package errors

import "fmt"

type NotFoundError struct {
	Entity string
	Field  string
	Value  string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s not found (%s = %s)", e.Entity, e.Field, e.Value)
}

type ConflictError struct {
	Field string
	Value string
}

func (e *ConflictError) Error() string {
	return fmt.Sprintf("%s already exists (%s)", e.Field, e.Value)
}

type ValidationError struct {
	Fields map[string]string
}

func (e *ValidationError) Error() string {
	return "validation error"
}

type UnauthorizedError struct {
	Reason string
}

func (e *UnauthorizedError) Error() string {
	return fmt.Sprintf("unauthorized: %s", e.Reason)
}
