package cart

import (
	"LojaGin/internal/modules/product"
	"LojaGin/internal/modules/user"

	"gorm.io/gorm"
)

// Cart representa o carrinho de compras de um usuário
type Cart struct {
	gorm.Model
	UserID uint       `json:"user_id" gorm:"unique;not null"`
	User   user.User  `json:"user"`
	Items  []CartItem `json:"items"`
}

// CartItem representa um item dentro do carrinho
type CartItem struct {
	gorm.Model
	CartID    uint            `json:"cart_id"`
	ProductID uint            `json:"product_id"`
	Product   product.Product `json:"product"`
	Quantity  uint            `json:"quantity" gorm:"not null;default:1"`
}

// Order representa um pedido feito após o checkout
type Order struct {
	gorm.Model
	UserID     uint        `json:"user_id"`
	User       user.User   `json:"user"`
	OrderItems []OrderItem `json:"order_items"`
	Total      float64     `json:"total"`
}

// OrderItem representa um item dentro de um pedido
type OrderItem struct {
	gorm.Model
	OrderID   uint            `json:"order_id"`
	ProductID uint            `json:"product_id"`
	Product   product.Product `json:"product"`
	Quantity  uint            `json:"quantity"`
	Price     float64         `json:"price"` // Preço no momento da compra
}

// AddToCartRequest define a estrutura para adicionar um item ao carrinho
type AddToCartRequest struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  uint `json:"quantity" binding:"required,min=1"`
}

// RemoveFromCartRequest define a estrutura para remover um item do carrinho
type RemoveFromCartRequest struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  uint `json:"quantity" binding:"required,min=1"`
}
