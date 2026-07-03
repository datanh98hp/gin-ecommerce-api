package validator

import "fmt"

// PaginationValidator provides validation for pagination parameters

// ValidatePaginationParams validates pagination parameters
func ValidatePaginationParams(page, pageSize int) error {
	if page < 1 {
		return fmt.Errorf("page must be greater than 0")
	}

	if page > 10000 {
		return fmt.Errorf("page must not exceed 10000")
	}

	if pageSize < 1 {
		return fmt.Errorf("page size must be greater than 0")
	}

	if pageSize > 100 {
		return fmt.Errorf("page size must not exceed 100")
	}

	return nil
}

// ValidateSortOrder validates sort order parameter
func ValidateSortOrder(sortOrder string) error {
	validOrders := map[string]bool{
		"asc":  true,
		"desc": true,
		"ASC":  true,
		"DESC": true,
	}

	if len(sortOrder) > 0 && !validOrders[sortOrder] {
		return fmt.Errorf("invalid sort order: must be 'asc' or 'desc'")
	}

	return nil
}

// ValidateSortBy validates sort by field parameter
func ValidateSortBy(sortBy string, allowedFields map[string]bool) error {
	if len(sortBy) == 0 {
		return nil // Optional
	}

	if !allowedFields[sortBy] {
		return fmt.Errorf("invalid sort field: %s", sortBy)
	}

	return nil
}

// ValidateDateRange validates date range parameters
func ValidateDateRange(startDate, endDate string) error {
	if len(startDate) == 0 && len(endDate) == 0 {
		return nil // Both optional
	}

	if len(startDate) > 0 && len(startDate) != 10 {
		return fmt.Errorf("invalid start date format: expected YYYY-MM-DD")
	}

	if len(endDate) > 0 && len(endDate) != 10 {
		return fmt.Errorf("invalid end date format: expected YYYY-MM-DD")
	}

	return nil
}

// PaginationParams holds validated pagination parameters
type PaginationParams struct {
	Page      int
	PageSize  int
	SortBy    string
	SortOrder string
}

// ValidateAndNormalizePaginationParams validates and normalizes pagination parameters with defaults
func ValidateAndNormalizePaginationParams(page, pageSize int, sortBy, sortOrder string, allowedSortFields map[string]bool) (*PaginationParams, error) {
	// Set defaults
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if len(sortOrder) == 0 {
		sortOrder = "desc"
	}

	// Validate
	if err := ValidatePaginationParams(page, pageSize); err != nil {
		return nil, err
	}

	if err := ValidateSortOrder(sortOrder); err != nil {
		return nil, err
	}

	if len(sortBy) > 0 && len(allowedSortFields) > 0 {
		if err := ValidateSortBy(sortBy, allowedSortFields); err != nil {
			return nil, err
		}
	}

	return &PaginationParams{
		Page:      page,
		PageSize:  pageSize,
		SortBy:    sortBy,
		SortOrder: sortOrder,
	}, nil
}

// CalculateOffset calculates database offset from page and page size
func CalculateOffset(page, pageSize int) int {
	return (page - 1) * pageSize
}

// GetLimit returns the limit for pagination
func GetLimit(pageSize int) int {
	return pageSize
}
