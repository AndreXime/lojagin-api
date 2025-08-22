package user

import "gorm.io/gorm"

type User struct {
	gorm.Model        // Fornece ID, CreatedAt, UpdatedAt, DeletedAt
	Name       string `json:"name"`
	Email      string `json:"email" gorm:"unique;not null"`
	Password   string `json:"-" gorm:"not null"`
}

type CreateUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type UpdateUserRequest struct {
	Name     *string `json:"name"`
	Email    *string `json:"email" binding:"omitempty,email"`
	Password *string `json:"password"`
}

type LoginUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
