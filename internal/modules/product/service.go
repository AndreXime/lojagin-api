package product

type Service struct {
	repo *Repository
}

func NewService(r *Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) Create(req CreateProductRequest) (*Product, error) {
	product := &Product{
		Name:       req.Name,
		Price:      req.Price,
		CategoryID: req.CategoryID,
	}
	return s.repo.CreateProduct(product)
}

func (s *Service) FindByID(id uint) (*Product, error) {
	return s.repo.GetProductByID(id)
}

func (s *Service) FindAll() ([]Product, error) {
	return s.repo.GetAllProducts()
}

func (s *Service) Update(id uint, req UpdateProductRequest) (*Product, error) {
	productToUpdate, err := s.repo.GetProductByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		productToUpdate.Name = *req.Name
	}
	if req.Price != nil {
		productToUpdate.Price = *req.Price
	}
	if req.CategoryID != nil {
		productToUpdate.CategoryID = *req.CategoryID
	}

	if err := s.repo.UpdateProduct(productToUpdate); err != nil {
		return nil, err
	}
	return productToUpdate, nil
}

func (s *Service) Delete(id uint) error {
	return s.repo.DeleteProduct(id)
}
