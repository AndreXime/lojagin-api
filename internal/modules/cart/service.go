package cart

type Service struct {
	repo *Repository
}

func NewService(r *Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) GetCart(userID uint) (*Cart, error) {
	return s.repo.GetOrCreateCart(userID)
}

func (s *Service) AddToCart(userID uint, req AddToCartRequest) (*Cart, error) {
	cart, err := s.repo.GetOrCreateCart(userID)
	if err != nil {
		return nil, err
	}
	if err := s.repo.AddItemToCart(cart.ID, req.ProductID, req.Quantity); err != nil {
		return nil, err
	}
	// Recarrega o carrinho para retornar o estado atualizado
	return s.repo.GetOrCreateCart(userID)
}

func (s *Service) RemoveFromCart(userID uint, req RemoveFromCartRequest) (*Cart, error) {
	cart, err := s.repo.GetOrCreateCart(userID)
	if err != nil {
		return nil, err
	}
	if err := s.repo.RemoveItemFromCart(cart.ID, req.ProductID, req.Quantity); err != nil {
		return nil, err
	}
	return s.repo.GetOrCreateCart(userID)
}

func (s *Service) ClearCart(userID uint) error {
	cart, err := s.repo.GetOrCreateCart(userID)
	if err != nil {
		return err
	}
	return s.repo.ClearCart(cart.ID)
}

func (s *Service) Checkout(userID uint) (*Order, error) {
	cart, err := s.repo.GetOrCreateCart(userID)
	if err != nil {
		return nil, err
	}
	if len(cart.Items) == 0 {
		return nil, ErrCartNotFound // Ou um erro "carrinho vazio"
	}
	return s.repo.Checkout(userID, cart)
}
