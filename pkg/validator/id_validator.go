package validator

import "fmt"

// IDValidator provides validation for ID parameters

// ValidateUserID validates user ID
func ValidateUserID(userID uint) error {
	if userID == 0 {
		return fmt.Errorf("user ID is required and must be greater than 0")
	}
	return nil
}

// ValidateIDList validates a list of IDs
func ValidateIDList(ids []uint) error {
	if len(ids) == 0 {
		return fmt.Errorf("ID list cannot be empty")
	}

	if len(ids) > 1000 {
		return fmt.Errorf("cannot process more than 1000 IDs at once")
	}

	for i, id := range ids {
		if id == 0 {
			return fmt.Errorf("ID at position %d is invalid (must be greater than 0)", i)
		}
	}

	return nil
}

// ValidateID validates a single generic ID
func ValidateID(id uint, fieldName string) error {
	if id == 0 {
		return fmt.Errorf("%s is required and must be greater than 0", fieldName)
	}
	return nil
}
