package category

type Service struct {
	repo *Repository
}

func NewService(r *Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) Create(req CreateCategoryRequest) (*Category, error) {
	category := &Category{
		Name: req.Name,
	}
	return s.repo.CreateCategory(category)
}

func (s *Service) FindByID(id uint) (*Category, error) {
	return s.repo.GetCategoryByID(id)
}

func (s *Service) FindAll() ([]Category, error) {
	return s.repo.GetAllCategories()
}

func (s *Service) Update(id uint, req UpdateCategoryRequest) (*Category, error) {
	categoryToUpdate, err := s.repo.GetCategoryByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		categoryToUpdate.Name = *req.Name
	}

	if err := s.repo.UpdateCategory(categoryToUpdate); err != nil {
		return nil, err
	}
	return categoryToUpdate, nil
}

func (s *Service) Delete(id uint) error {
	return s.repo.DeleteCategory(id)
}
