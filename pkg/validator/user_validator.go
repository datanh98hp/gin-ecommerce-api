package validator

import "fmt"

// UserValidator provides validation for user-related requests

// ValidateLoginRequest validates login request
func ValidateLoginRequest(email, password string) error {
	if len(email) == 0 {
		return fmt.Errorf("email is required")
	}

	if !IsValidEmail(email) {
		return fmt.Errorf("invalid email format")
	}

	if len(password) == 0 {
		return fmt.Errorf("password is required")
	}

	if len(password) < 6 {
		return fmt.Errorf("password must be at least 6 characters")
	}

	if len(password) > 128 {
		return fmt.Errorf("password must not exceed 128 characters")
	}

	return nil
}

// ValidateRegisterRequest validates registration request
func ValidateRegisterRequest(email, username, password, firstName, lastName string) error {
	// Validate email
	if len(email) == 0 {
		return fmt.Errorf("email is required")
	}

	if !IsValidEmail(email) {
		return fmt.Errorf("invalid email format")
	}

	// Validate username
	if len(username) == 0 {
		return fmt.Errorf("username is required")
	}

	if !IsValidUsername(username) {
		return fmt.Errorf("username must be 3-30 characters long and contain only alphanumeric characters, underscores, or hyphens")
	}

	// Validate password
	if len(password) == 0 {
		return fmt.Errorf("password is required")
	}

	if !IsValidPassword(password) {
		return fmt.Errorf("password must be 6-128 characters long")
	}

	// Validate first name
	if len(firstName) == 0 {
		return fmt.Errorf("first name is required")
	}

	if !IsValidName(firstName) {
		return fmt.Errorf("first name must be 2-50 characters long and contain only letters, spaces, hyphens, or apostrophes")
	}

	// Validate last name
	if len(lastName) == 0 {
		return fmt.Errorf("last name is required")
	}

	if !IsValidName(lastName) {
		return fmt.Errorf("last name must be 2-50 characters long and contain only letters, spaces, hyphens, or apostrophes")
	}

	return nil
}

// ValidateUpdateUserRequest validates user update request
func ValidateUpdateUserRequest(firstName, lastName *string, role *string) error {
	if firstName != nil {
		if len(*firstName) == 0 {
			return fmt.Errorf("first name cannot be empty")
		}
		if !IsValidName(*firstName) {
			return fmt.Errorf("first name must be 2-50 characters long and contain only letters, spaces, hyphens, or apostrophes")
		}
	}

	if lastName != nil {
		if len(*lastName) == 0 {
			return fmt.Errorf("last name cannot be empty")
		}
		if !IsValidName(*lastName) {
			return fmt.Errorf("last name must be 2-50 characters long and contain only letters, spaces, hyphens, or apostrophes")
		}
	}

	if role != nil {
		if !IsValidUserRole(*role) {
			return fmt.Errorf("invalid role: must be 'user' or 'admin'")
		}
	}

	return nil
}

// ValidateChangePasswordRequest validates change password request
func ValidateChangePasswordRequest(oldPassword, newPassword, confirmPassword string) error {
	if len(oldPassword) == 0 {
		return fmt.Errorf("old password is required")
	}

	if len(newPassword) == 0 {
		return fmt.Errorf("new password is required")
	}

	if !IsValidPassword(newPassword) {
		return fmt.Errorf("new password must be 6-128 characters long")
	}

	if len(confirmPassword) == 0 {
		return fmt.Errorf("password confirmation is required")
	}

	if newPassword != confirmPassword {
		return fmt.Errorf("passwords do not match")
	}

	if oldPassword == newPassword {
		return fmt.Errorf("new password must be different from old password")
	}

	return nil
}
