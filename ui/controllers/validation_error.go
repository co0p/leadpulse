package controllers

// ValidationError represents a validation failure with a user-friendly message.
// Low-level errors are caught and wrapped as ValidationError with clear messaging.
type ValidationError struct {
	// Kind categorizes the error for the UI to handle appropriately
	Kind ValidationErrorKind
	// Message is a user-friendly error message suitable for display
	Message string
	// Internal error is captured but not exposed to the UI
	internal error
}

type ValidationErrorKind string

const (
	// ValidationErrorKindNoMemberSelected — user must select a member before saving
	ValidationErrorKindNoMemberSelected ValidationErrorKind = "no_member_selected"

	// ValidationErrorKindNoData — no data entered to save
	ValidationErrorKindNoData ValidationErrorKind = "no_data"

	// ValidationErrorKindInvalidField — a specific field contains invalid data
	ValidationErrorKindInvalidField ValidationErrorKind = "invalid_field"

	// ValidationErrorKindDatabaseFailure — database operation failed (low-level error wrapped)
	ValidationErrorKindDatabaseFailure ValidationErrorKind = "database_failure"

	// ValidationErrorKindUnknown — unexpected error; should not happen in normal usage
	ValidationErrorKindUnknown ValidationErrorKind = "unknown"
)

func (e *ValidationError) Error() string {
	return e.Message
}

// IsValidationError checks if an error is a ValidationError.
func IsValidationError(err error) bool {
	_, ok := err.(*ValidationError)
	return ok
}

// GetValidationErrorKind returns the kind of a ValidationError, or an empty string if not a ValidationError.
func GetValidationErrorKind(err error) ValidationErrorKind {
	if ve, ok := err.(*ValidationError); ok {
		return ve.Kind
	}
	return ""
}

// NewValidationError creates a new ValidationError with the given kind and user-facing message.
// The internal error is captured but never exposed to the UI.
func NewValidationError(kind ValidationErrorKind, message string, internal error) *ValidationError {
	return &ValidationError{
		Kind:     kind,
		Message:  message,
		internal: internal,
	}
}

// WrapError wraps a low-level error with a user-friendly message and kind.
func WrapError(kind ValidationErrorKind, message string, err error) *ValidationError {
	return &ValidationError{
		Kind:     kind,
		Message:  message,
		internal: err,
	}
}

// UserMessage returns the user-facing message. Never exposes internal details.
func (e *ValidationError) UserMessage() string {
	return e.Message
}
