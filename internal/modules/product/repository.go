package product

import (
	"errors"

	"gorm.io/gorm"
)

var (
	ErrProductNotFound = errors.New("produto não encontrado")
	ErrDatabase        = errors.New("ocorreu um erro no banco de dados")
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateProduct(product *Product) (*Product, error) {
	if err := r.db.Create(product).Error; err != nil {
		return nil, ErrDatabase
	}
	return product, nil
}

func (r *Repository) GetProductByID(id uint) (*Product, error) {
	var product Product
	if err := r.db.Preload("Category").First(&product, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, ErrDatabase
	}
	return &product, nil
}

func (r *Repository) GetAllProducts() ([]Product, error) {
	var products []Product
	if err := r.db.Preload("Category").Find(&products).Error; err != nil {
		return nil, ErrDatabase
	}
	return products, nil
}

func (r *Repository) UpdateProduct(product *Product) error {
	if err := r.db.Save(product).Error; err != nil {
		return ErrDatabase
	}
	return nil
}

func (r *Repository) DeleteProduct(id uint) error {
	result := r.db.Delete(&Product{}, id)
	if result.Error != nil {
		return ErrDatabase
	}
	if result.RowsAffected == 0 {
		return ErrProductNotFound
	}
	return nil
}
