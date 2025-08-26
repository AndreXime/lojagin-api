package cart

import (
	"errors"

	"gorm.io/gorm"
)

var (
	ErrCartNotFound      = errors.New("carrinho não encontrado")
	ErrProductNotInCart  = errors.New("produto não encontrado no carrinho")
	ErrInsufficientStock = errors.New("a quantidade a remover é maior que a existente no carrinho")
	ErrDatabase          = errors.New("ocorreu um erro no banco de dados")
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// GetOrCreateCart busca um carrinho para o usuário ou cria um novo se não existir
func (r *Repository) GetOrCreateCart(userID uint) (*Cart, error) {
	var cart Cart
	err := r.db.Preload("Items.Product.Category").Where(Cart{UserID: userID}).FirstOrCreate(&cart).Error
	if err != nil {
		return nil, ErrDatabase
	}
	return &cart, nil
}

// AddItemToCart adiciona ou atualiza um item no carrinho
func (r *Repository) AddItemToCart(cartID, productID, quantity uint) error {
	var cartItem CartItem
	err := r.db.Where(CartItem{CartID: cartID, ProductID: productID}).First(&cartItem).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrDatabase
	}

	if errors.Is(err, gorm.ErrRecordNotFound) { // Novo item
		cartItem = CartItem{
			CartID:    cartID,
			ProductID: productID,
			Quantity:  quantity,
		}
		if err := r.db.Create(&cartItem).Error; err != nil {
			return ErrDatabase
		}
	} else { // Item existente, atualiza a quantidade
		cartItem.Quantity += quantity
		if err := r.db.Save(&cartItem).Error; err != nil {
			return ErrDatabase
		}
	}
	return nil
}

// RemoveItemFromCart remove ou diminui a quantidade de um item no carrinho
func (r *Repository) RemoveItemFromCart(cartID, productID, quantity uint) error {
	var cartItem CartItem
	if err := r.db.Where(CartItem{CartID: cartID, ProductID: productID}).First(&cartItem).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrProductNotInCart
		}
		return ErrDatabase
	}

	if cartItem.Quantity < quantity {
		return ErrInsufficientStock
	}

	if cartItem.Quantity == quantity { // Remove completamente o item
		if err := r.db.Delete(&cartItem).Error; err != nil {
			return ErrDatabase
		}
	} else { // Apenas diminui a quantidade
		cartItem.Quantity -= quantity
		if err := r.db.Save(&cartItem).Error; err != nil {
			return ErrDatabase
		}
	}
	return nil
}

// ClearCart remove todos os itens de um carrinho
func (r *Repository) ClearCart(cartID uint) error {
	if err := r.db.Where("cart_id = ?", cartID).Delete(&CartItem{}).Error; err != nil {
		return ErrDatabase
	}
	return nil
}

// Checkout cria um pedido a partir do carrinho e o limpa
func (r *Repository) Checkout(userID uint, cart *Cart) (*Order, error) {
	order := &Order{
		UserID:     userID,
		OrderItems: []OrderItem{},
		Total:      0,
	}

	// Transação para garantir consistência
	err := r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Cria o Pedido
		if err := tx.Create(order).Error; err != nil {
			return err
		}

		var total float64 = 0
		// 2. Move os itens do carrinho para itens de pedido
		for _, item := range cart.Items {
			orderItem := OrderItem{
				OrderID:   order.ID,
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				Price:     item.Product.Price, // Salva o preço atual do produto
			}
			if err := tx.Create(&orderItem).Error; err != nil {
				return err // Rollback
			}
			total += item.Product.Price * float64(item.Quantity)
		}

		// 3. Atualiza o total do pedido
		order.Total = total
		if err := tx.Save(order).Error; err != nil {
			return err
		}

		// 4. Limpa o carrinho
		if err := tx.Where("cart_id = ?", cart.ID).Delete(CartItem{}).Error; err != nil {
			return err // Rollback
		}

		return nil
	})

	if err != nil {
		return nil, ErrDatabase
	}

	// Recarrega o pedido com os itens para retornar a resposta completa
	r.db.Preload("OrderItems.Product.Category").First(order, order.ID)

	return order, nil
}
