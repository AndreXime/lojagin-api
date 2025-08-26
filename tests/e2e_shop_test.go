package tests

import (
	"LojaGin/internal/modules/category"
	"LojaGin/internal/modules/product"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// setupCategory é uma função helper para criar uma categoria e retornar seu ID
func setupCategory(t *testing.T, authToken string) uint {
	categoryPayload := []byte(`{"name": "Eletrônicos"}`)
	req, _ := http.NewRequest(http.MethodPost, "/api/categories/", bytes.NewBuffer(categoryPayload))
	req.Header.Set("Authorization", "Bearer "+authToken)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var createdCategory category.Category
	err := json.Unmarshal(w.Body.Bytes(), &createdCategory)
	assert.NoError(t, err)
	assert.NotZero(t, createdCategory.ID)

	return createdCategory.ID
}

func TestCategoryEndpoints(t *testing.T) {
	clearDB()
	_, authToken := setupAuthenticatedUser(t)

	// --- Cenário 1: Tentar criar categoria sem autenticação ---
	t.Run("POST /categories - Unauthorized", func(t *testing.T) {
		categoryPayload := []byte(`{"name": "Livros"}`)
		req, _ := http.NewRequest(http.MethodPost, "/api/categories/", bytes.NewBuffer(categoryPayload))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	// --- Cenário 2: Criar categoria com sucesso ---
	t.Run("POST /categories - Success with Auth", func(t *testing.T) {
		setupCategory(t, authToken)
	})

	// --- Cenário 3: Listar categorias (público) ---
	t.Run("GET /categories - Public Success", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/categories/", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		var categories []category.Category
		json.Unmarshal(w.Body.Bytes(), &categories)
		assert.NotEmpty(t, categories)
	})

	// --- Cenário 4: Atualizar uma categoria com sucesso ---
	t.Run("PUT /categories/:id - Success with Auth", func(t *testing.T) {
		categoryID := setupCategory(t, authToken)
		updatePayload := []byte(`{"name": "Eletrônicos Modernos"}`)
		url := fmt.Sprintf("/api/categories/%d", categoryID)
		req, _ := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(updatePayload))
		req.Header.Set("Authorization", "Bearer "+authToken)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		var updatedCategory category.Category
		json.Unmarshal(w.Body.Bytes(), &updatedCategory)
		assert.Equal(t, "Eletrônicos Modernos", updatedCategory.Name)
	})

	// --- Cenário 5: Deletar uma categoria com sucesso ---
	t.Run("DELETE /categories/:id - Success with Auth", func(t *testing.T) {
		categoryID := setupCategory(t, authToken)
		url := fmt.Sprintf("/api/categories/%d", categoryID)
		req, _ := http.NewRequest(http.MethodDelete, url, nil)
		req.Header.Set("Authorization", "Bearer "+authToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNoContent, w.Code)
	})
}

func TestProductEndpoints(t *testing.T) {
	clearDB()
	_, authToken := setupAuthenticatedUser(t)
	categoryID := setupCategory(t, authToken)

	// --- Cenário 1: Tentar criar produto sem autenticação ---
	t.Run("POST /products - Unauthorized", func(t *testing.T) {
		productPayload := []byte(fmt.Sprintf(`{"name": "Notebook", "price": 4500.00, "category_id": %d}`, categoryID))
		req, _ := http.NewRequest(http.MethodPost, "/api/products/", bytes.NewBuffer(productPayload))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	// --- Cenário 2: Criar produto com sucesso ---
	var productID uint
	t.Run("POST /products - Success with Auth", func(t *testing.T) {
		productPayload := []byte(fmt.Sprintf(`{"name": "Notebook Pro", "price": 5500.00, "category_id": %d}`, categoryID))
		req, _ := http.NewRequest(http.MethodPost, "/api/products/", bytes.NewBuffer(productPayload))
		req.Header.Set("Authorization", "Bearer "+authToken)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)

		var createdProduct product.Product
		json.Unmarshal(w.Body.Bytes(), &createdProduct)
		assert.Equal(t, "Notebook Pro", createdProduct.Name)
		assert.Equal(t, categoryID, createdProduct.CategoryID)
		productID = createdProduct.ID
	})

	// --- Cenário 3: Listar produtos (público) ---
	t.Run("GET /products - Public Success", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/products/", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		var products []product.Product
		json.Unmarshal(w.Body.Bytes(), &products)
		assert.NotEmpty(t, products)
	})

	// --- Cenário 4: Obter produto por ID (público) ---
	t.Run("GET /products/:id - Public Success", func(t *testing.T) {
		url := fmt.Sprintf("/api/products/%d", productID)
		req, _ := http.NewRequest(http.MethodGet, url, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		var foundProduct product.Product
		json.Unmarshal(w.Body.Bytes(), &foundProduct)
		assert.Equal(t, productID, foundProduct.ID)
	})

	// --- Cenário 5: Atualizar um produto com sucesso ---
	t.Run("PUT /products/:id - Success with Auth", func(t *testing.T) {
		newPrice := 5250.50
		updatePayload := []byte(fmt.Sprintf(`{"price": %f}`, newPrice))
		url := fmt.Sprintf("/api/products/%d", productID)
		req, _ := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(updatePayload))
		req.Header.Set("Authorization", "Bearer "+authToken)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		var updatedProduct product.Product
		json.Unmarshal(w.Body.Bytes(), &updatedProduct)
		assert.Equal(t, newPrice, updatedProduct.Price)
	})

	// --- Cenário 6: Deletar um produto com sucesso ---
	t.Run("DELETE /products/:id - Success with Auth", func(t *testing.T) {
		url := fmt.Sprintf("/api/products/%d", productID)
		req, _ := http.NewRequest(http.MethodDelete, url, nil)
		req.Header.Set("Authorization", "Bearer "+authToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNoContent, w.Code)
	})
}
