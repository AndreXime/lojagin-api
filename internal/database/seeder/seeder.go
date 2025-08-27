package seeder

import (
	"LojaGin/internal/modules/category"
	"LojaGin/internal/modules/product"
	"LojaGin/internal/modules/user"
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Run executa os seeders para popular o banco de dados com dados iniciais.
func Run(db *gorm.DB) {
	if err := seedUsers(db); err != nil {
		log.Fatalf("Não foi possível popular usuários: %v", err)
	}

	categories, err := seedCategories(db)
	if err != nil {
		log.Fatalf("Não foi possível popular categorias: %v", err)
	}

	if err := seedProducts(db, categories); err != nil {
		log.Fatalf("Não foi possível popular produtos: %v", err)
	}

	log.Println("Seeder executado com sucesso!")
}

// seedUsers insere usuários iniciais se a tabela estiver vazia.
func seedUsers(db *gorm.DB) error {
	var count int64
	db.Model(&user.User{}).Count(&count)
	if count > 0 {
		return nil // Usuários já existem
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("12345678"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	users := []user.User{
		{Name: "Admin User", Email: "admin@example.com", Password: string(hashedPassword)},
		{Name: "Test User", Email: "test@example.com", Password: string(hashedPassword)},
	}

	return db.Create(&users).Error
}

// seedCategories insere categorias iniciais e as retorna.
func seedCategories(db *gorm.DB) (map[string]uint, error) {
	var count int64
	db.Model(&category.Category{}).Count(&count)
	if count > 0 {
		// Se já existem, apenas as busca para retornar os IDs
		var existingCategories []category.Category
		db.Find(&existingCategories)
		categoryMap := make(map[string]uint)
		for _, cat := range existingCategories {
			categoryMap[cat.Name] = cat.ID
		}
		return categoryMap, nil
	}

	categories := []category.Category{
		{Name: "Eletrônicos"},
		{Name: "Livros"},
		{Name: "Casa e Cozinha"},
	}

	if err := db.Create(&categories).Error; err != nil {
		return nil, err
	}

	categoryMap := make(map[string]uint)
	for _, cat := range categories {
		categoryMap[cat.Name] = cat.ID
	}

	return categoryMap, nil
}

// seedProducts insere produtos iniciais associados às categorias.
func seedProducts(db *gorm.DB, categories map[string]uint) error {
	var count int64
	db.Model(&product.Product{}).Count(&count)
	if count > 0 {
		return nil // Produtos já existem
	}

	products := []product.Product{
		{Name: "Notebook Gamer Pro", Price: 7500.99, CategoryID: categories["Eletrônicos"]},
		{Name: "Mouse Vertical Sem Fio", Price: 189.90, CategoryID: categories["Eletrônicos"]},
		{Name: "O Guia do Mochileiro das Galáxias", Price: 42.00, CategoryID: categories["Livros"]},
		{Name: "Fritadeira Elétrica Sem Óleo", Price: 399.50, CategoryID: categories["Casa e Cozinha"]},
		{Name: "Jogo de Panelas Antiaderente", Price: 250.00, CategoryID: categories["Casa e Cozinha"]},
	}

	return db.Create(&products).Error
}
