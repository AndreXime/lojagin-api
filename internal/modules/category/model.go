package category

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Name string `json:"name" gorm:"unique;not null"`
}

type CreateCategoryRequest struct {
	Name string `json:"name" binding:"required"`
}

type UpdateCategoryRequest struct {
	Name *string `json:"name"`
}
