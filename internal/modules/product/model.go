package product

import (
	"LojaGin/internal/modules/category"

	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name       string            `json:"name" gorm:"not null"`
	Price      float64           `json:"price" gorm:"not null"`
	CategoryID uint              `json:"category_id"`
	Category   category.Category `json:"category"`
}

type CreateProductRequest struct {
	Name       string  `json:"name" binding:"required"`
	Price      float64 `json:"price" binding:"required,gt=0"`
	CategoryID uint    `json:"category_id" binding:"required"`
}

type UpdateProductRequest struct {
	Name       *string  `json:"name"`
	Price      *float64 `json:"price"`
	CategoryID *uint    `json:"category_id"`
}
