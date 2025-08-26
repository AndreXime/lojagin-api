package category

import (
	"errors"

	"gorm.io/gorm"
)

var (
	ErrCategoryNotFound = errors.New("categoria não encontrada")
	ErrDatabase         = errors.New("ocorreu um erro no banco de dados")
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateCategory(category *Category) (*Category, error) {
	if err := r.db.Create(category).Error; err != nil {
		return nil, ErrDatabase
	}
	return category, nil
}

func (r *Repository) GetCategoryByID(id uint) (*Category, error) {
	var category Category
	if err := r.db.First(&category, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCategoryNotFound
		}
		return nil, ErrDatabase
	}
	return &category, nil
}

func (r *Repository) GetAllCategories() ([]Category, error) {
	var categories []Category
	if err := r.db.Find(&categories).Error; err != nil {
		return nil, ErrDatabase
	}
	return categories, nil
}

func (r *Repository) UpdateCategory(category *Category) error {
	if err := r.db.Save(category).Error; err != nil {
		return ErrDatabase
	}
	return nil
}

func (r *Repository) DeleteCategory(id uint) error {
	result := r.db.Delete(&Category{}, id)
	if result.Error != nil {
		return ErrDatabase
	}
	if result.RowsAffected == 0 {
		return ErrCategoryNotFound
	}
	return nil
}
