package validator

// Validator is the main package for API request validation
// It provides a comprehensive set of validators for all API endpoints

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

// ValidationErrors is a collection of validation errors
type ValidationErrors []ValidationError

// Error implements the error interface
func (ve ValidationErrors) Error() string {
	if len(ve) == 0 {
		return "validation passed"
	}
	return ve[0].Message
}

// Add adds a validation error
func (ve *ValidationErrors) Add(field, message string) {
	*ve = append(*ve, ValidationError{Field: field, Message: message})
}

// HasErrors checks if there are any validation errors
func (ve ValidationErrors) HasErrors() bool {
	return len(ve) > 0
}

// All returns all validation errors
func (ve ValidationErrors) All() []ValidationError {
	return ve
}

// First returns the first validation error
func (ve ValidationErrors) First() *ValidationError {
	if len(ve) > 0 {
		return &ve[0]
	}
	return nil
}

// ByField returns all errors for a specific field
func (ve ValidationErrors) ByField(field string) []ValidationError {
	var errors []ValidationError
	for _, err := range ve {
		if err.Field == field {
			errors = append(errors, err)
		}
	}
	return errors
}

// NewValidationErrors creates a new ValidationErrors instance
func NewValidationErrors() ValidationErrors {
	return ValidationErrors{}
}

// ValidateStruct provides a generic validation interface for custom validation
type ValidateStruct interface {
	Validate() error
}
