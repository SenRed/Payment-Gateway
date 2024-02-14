package model

type ErrorType string

const (
	DuplicatedSessionId ErrorType = "DuplicatedSessionId"
	SessionNotFound     ErrorType = "SessionNotFound"
)

type ErrorCategory string

const (
	FunctionalError ErrorCategory = "FunctionalError"
	TechnicalError  ErrorCategory = "TechnicalError"
)

type DomainError struct {
	Category  ErrorCategory
	Type      ErrorType
	Message   string
	RootCause error
}
