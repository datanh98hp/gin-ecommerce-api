package validator

import "fmt"

// PostValidator provides validation for post-related requests

// ValidateCreatePostRequest validates post creation request
func ValidateCreatePostRequest(title, content string) error {
	// Validate title
	if len(title) == 0 {
		return fmt.Errorf("post title is required")
	}

	title = TrimWhitespace(title)
	if len(title) < 5 {
		return fmt.Errorf("post title must be at least 5 characters long")
	}

	if len(title) > 255 {
		return fmt.Errorf("post title must not exceed 255 characters")
	}

	// Validate content
	if len(content) == 0 {
		return fmt.Errorf("post content is required")
	}

	content = TrimWhitespace(content)
	if len(content) < 20 {
		return fmt.Errorf("post content must be at least 20 characters long")
	}

	if len(content) > 10000 {
		return fmt.Errorf("post content must not exceed 10000 characters")
	}

	return nil
}

// ValidateUpdatePostRequest validates post update request
func ValidateUpdatePostRequest(title *string, content *string, isActive *bool) error {
	if title != nil {
		if len(*title) == 0 {
			return fmt.Errorf("post title cannot be empty")
		}
		trimmedTitle := TrimWhitespace(*title)
		if len(trimmedTitle) < 5 {
			return fmt.Errorf("post title must be at least 5 characters long")
		}
		if len(trimmedTitle) > 255 {
			return fmt.Errorf("post title must not exceed 255 characters")
		}
	}

	if content != nil {
		if len(*content) == 0 {
			return fmt.Errorf("post content cannot be empty")
		}
		trimmedContent := TrimWhitespace(*content)
		if len(trimmedContent) < 20 {
			return fmt.Errorf("post content must be at least 20 characters long")
		}
		if len(trimmedContent) > 10000 {
			return fmt.Errorf("post content must not exceed 10000 characters")
		}
	}

	// isActive is optional and doesn't need validation if it's a boolean

	return nil
}

// ValidatePostID validates post ID
func ValidatePostID(postID uint) error {
	if postID == 0 {
		return fmt.Errorf("post ID is required and must be greater than 0")
	}
	return nil
}

// ValidatePostSearch validates post search parameters
func ValidatePostSearch(searchQuery string) error {
	if len(searchQuery) > 255 {
		return fmt.Errorf("search query must not exceed 255 characters")
	}
	return nil
}
