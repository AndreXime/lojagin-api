package tests

import (
	"LojaGin/internal/server"
	"LojaGin/internal/user"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestUserEndpoints(t *testing.T) {
	router := gin.Default()
	server.SetupAPI(router)

	// Variáveis para compartilhar estado entre os sub-testes
	var authToken string
	var createdUser user.User

	// --- Cenário 1: Registrar um novo usuário com sucesso ---
	t.Run("should register user successfully", func(t *testing.T) {
		// A rota de registro exige nome, email e senha
		registerPayload := []byte(`{"name": "Marie Curie", "email": "marie@radium.com", "password": "password123"}`)

		// A rota correta é /api/auth/register
		req, _ := http.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(registerPayload))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// A rota de registro retorna 200 OK com o usuário criado
		assert.Equal(t, http.StatusOK, w.Code)

		err := json.Unmarshal(w.Body.Bytes(), &createdUser)
		assert.NoError(t, err)
		assert.Equal(t, "Marie Curie", createdUser.Name)
		assert.Equal(t, "marie@radium.com", createdUser.Email)
		assert.Empty(t, createdUser.Password, "A senha não deve ser retornada na resposta")
		assert.NotZero(t, createdUser.ID)
	})

	// --- Cenário 2: Fazer login para obter o token JWT ---
	t.Run("should login and get JWT token", func(t *testing.T) {
		// Garante que o teste de registro foi executado e temos um usuário
		if createdUser.ID == 0 {
			t.Skip("Skipping login test because user registration failed")
		}

		loginPayload := []byte(`{"email": "marie@radium.com", "password": "password123"}`)

		// A rota correta é /api/auth/login
		req, _ := http.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(loginPayload))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var loginResponse map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &loginResponse)
		assert.NoError(t, err)
		assert.Contains(t, loginResponse, "token", "A resposta de login deve conter um token")

		authToken = loginResponse["token"]
		assert.NotEmpty(t, authToken)
	})

	// --- Cenário 3: Buscar um usuário pelo ID (requer autenticação) ---
	t.Run("should get user by id with auth", func(t *testing.T) {
		if authToken == "" {
			t.Skip("Skipping get user test because auth token was not retrieved")
		}

		// A rota correta é /api/users/:id
		url := fmt.Sprintf("/api/users/%d", createdUser.ID)
		req, _ := http.NewRequest(http.MethodGet, url, nil)
		// Adiciona o token de autenticação no cabeçalho
		req.Header.Set("Authorization", "Bearer "+authToken)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var foundUser user.User
		err := json.Unmarshal(w.Body.Bytes(), &foundUser)
		assert.NoError(t, err)
		assert.Equal(t, createdUser.Name, foundUser.Name)
		assert.Equal(t, createdUser.Email, foundUser.Email)
		assert.Equal(t, createdUser.ID, foundUser.ID)
	})

	// --- Cenário 4: Tentar buscar um usuário que não existe ---
	t.Run("should return not found for non-existent user", func(t *testing.T) {
		if authToken == "" {
			t.Skip("Skipping get user test because auth token was not retrieved")
		}

		req, _ := http.NewRequest(http.MethodGet, "/api/users/99999", nil)
		req.Header.Set("Authorization", "Bearer "+authToken)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// O status esperado é Not Found
		assert.Equal(t, http.StatusNotFound, w.Code)

		var errorResponse map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
		assert.NoError(t, err)
		// A mensagem de erro correta vem de `ErrUserNotFound`
		assert.Equal(t, "usuário não encontrado", errorResponse["error"])
	})

	// --- Cenário 5: Tentar acessar uma rota protegida sem token ---
	t.Run("should fail to get user without auth token", func(t *testing.T) {
		url := fmt.Sprintf("/api/users/%d", createdUser.ID)
		req, _ := http.NewRequest(http.MethodGet, url, nil)
		// NENHUM token é enviado

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// O middleware de autenticação deve barrar a requisição
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
