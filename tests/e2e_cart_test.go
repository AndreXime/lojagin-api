package tests

import (
	"LojaGin/internal/modules/cart"
	"LojaGin/internal/modules/product"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// setupProduct é uma função helper para criar um produto e retornar seu ID
func setupProduct(t *testing.T, authToken string, categoryID uint) uint {
	productPayload := []byte(fmt.Sprintf(`{"name": "Mouse Gamer", "price": 150.00, "category_id": %d}`, categoryID))
	req, _ := http.NewRequest(http.MethodPost, "/api/products/", bytes.NewBuffer(productPayload))
	req.Header.Set("Authorization", "Bearer "+authToken)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	var createdProduct product.Product
	err := json.Unmarshal(w.Body.Bytes(), &createdProduct)
	assert.NoError(t, err)
	return createdProduct.ID
}

func TestCartEndpoints(t *testing.T) {
	clearDB()
	_, authToken := setupAuthenticatedUser(t)
	categoryID := setupCategory(t, authToken)
	productID := setupProduct(t, authToken, categoryID)

	// --- Cenário 1: Obter carrinho vazio ---
	t.Run("GET /cart - Get Empty Cart", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/cart/", nil)
		req.Header.Set("Authorization", "Bearer "+authToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		var cartResponse cart.Cart
		json.Unmarshal(w.Body.Bytes(), &cartResponse)
		assert.Empty(t, cartResponse.Items)
	})

	// --- Cenário 2: Adicionar item ao carrinho ---
	t.Run("POST /cart/add - Add Item", func(t *testing.T) {
		addPayload := []byte(fmt.Sprintf(`{"product_id": %d, "quantity": 2}`, productID))
		req, _ := http.NewRequest(http.MethodPost, "/api/cart/add", bytes.NewBuffer(addPayload))
		req.Header.Set("Authorization", "Bearer "+authToken)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		var cartResponse cart.Cart
		json.Unmarshal(w.Body.Bytes(), &cartResponse)
		assert.Len(t, cartResponse.Items, 1)
		assert.Equal(t, uint(2), cartResponse.Items[0].Quantity)
		assert.Equal(t, productID, cartResponse.Items[0].ProductID)
	})

	// --- Cenário 3: Remover uma unidade do item do carrinho ---
	t.Run("POST /cart/remove - Remove Partial Quantity", func(t *testing.T) {
		removePayload := []byte(fmt.Sprintf(`{"product_id": %d, "quantity": 1}`, productID))
		req, _ := http.NewRequest(http.MethodPost, "/api/cart/remove", bytes.NewBuffer(removePayload))
		req.Header.Set("Authorization", "Bearer "+authToken)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		var cartResponse cart.Cart
		json.Unmarshal(w.Body.Bytes(), &cartResponse)
		assert.Len(t, cartResponse.Items, 1)
		assert.Equal(t, uint(1), cartResponse.Items[0].Quantity)
	})

	// --- Cenário 4: Limpar o carrinho ---
	t.Run("DELETE /cart/clear - Clear Cart", func(t *testing.T) {
		// Primeiro, adiciona um item para garantir que não está vazio
		addPayload := []byte(fmt.Sprintf(`{"product_id": %d, "quantity": 1}`, productID))
		req, _ := http.NewRequest(http.MethodPost, "/api/cart/add", bytes.NewBuffer(addPayload))
		req.Header.Set("Authorization", "Bearer "+authToken)
		router.ServeHTTP(httptest.NewRecorder(), req)

		// Agora, limpa
		req, _ = http.NewRequest(http.MethodDelete, "/api/cart/clear", nil)
		req.Header.Set("Authorization", "Bearer "+authToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNoContent, w.Code)

		// Verifica se o carrinho está vazio
		req, _ = http.NewRequest(http.MethodGet, "/api/cart/", nil)
		req.Header.Set("Authorization", "Bearer "+authToken)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		var cartResponse cart.Cart
		json.Unmarshal(w.Body.Bytes(), &cartResponse)
		assert.Empty(t, cartResponse.Items)
	})

	// --- Cenário 5: Checkout ---
	t.Run("POST /cart/checkout - Success", func(t *testing.T) {
		// Adiciona itens ao carrinho novamente para o checkout
		addPayload := []byte(fmt.Sprintf(`{"product_id": %d, "quantity": 3}`, productID))
		req, _ := http.NewRequest(http.MethodPost, "/api/cart/add", bytes.NewBuffer(addPayload))
		req.Header.Set("Authorization", "Bearer "+authToken)
		router.ServeHTTP(httptest.NewRecorder(), req)

		// Faz o checkout
		req, _ = http.NewRequest(http.MethodPost, "/api/cart/checkout", nil)
		req.Header.Set("Authorization", "Bearer "+authToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)

		var orderResponse cart.Order
		json.Unmarshal(w.Body.Bytes(), &orderResponse)
		assert.Len(t, orderResponse.OrderItems, 1)
		assert.Equal(t, float64(450), orderResponse.Total) // 3 * 150.00

		// Verifica se o carrinho foi esvaziado após o checkout
		req, _ = http.NewRequest(http.MethodGet, "/api/cart/", nil)
		req.Header.Set("Authorization", "Bearer "+authToken)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		var cartResponse cart.Cart
		json.Unmarshal(w.Body.Bytes(), &cartResponse)
		assert.Empty(t, cartResponse.Items)
	})
}
