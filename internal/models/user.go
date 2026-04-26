package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Email     string         `gorm:"uniqueIndex;not null" json:"email"`
	Username  string         `gorm:"uniqueIndex;not null" json:"username"`
	Password  string         `gorm:"not null" json:"password"`
	FirstName string         `json:"first_name"`
	LastName  string         `json:"last_name"`
	Role      string         `gorm:"default:'user'" json:"role"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Cart   *Cart   `gorm:"foreignKey:UserID" json:"cart,omitempty"`
	Orders []Order `gorm:"foreignKey:UserID" json:"orders,omitempty"`
	Posts  []Post  `gorm:"foreignKey:UserID" json:"posts,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Username  string `json:"username" binding:"required,min=3"`
	Password  string `json:"password" binding:"required,min=6"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

type UserResponse struct {
	Email     string         `gorm:"uniqueIndex;not null" json:"email"`
	Username  string         `gorm:"uniqueIndex;not null" json:"username"`
	FirstName string         `json:"first_name"`
	LastName  string         `json:"last_name"`
	Role      string         `gorm:"default:'user'" json:"role"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
}
type AuthResponse struct {
	Token string `json:"token"`
	User  UserResponse   `json:"user"`
}
type UpdateUserRequest struct {
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Email     *string `json:"email"`
	IsActive  *bool   `json:"is_active"`
}
