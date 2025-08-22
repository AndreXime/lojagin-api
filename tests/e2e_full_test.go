package tests

import (
	"LojaGin/internal/config"
	"LojaGin/internal/database"
	"LojaGin/internal/modules/user"
	"LojaGin/internal/routes"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var router *gin.Engine
var db *gorm.DB

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	// 2. Define um nome para o banco de dados de teste.
	config.DB_URL = "test_e2e.db"
	config.JWT_SECRET = "segredo_para_testes_123456"

	// 3. Garante que qualquer arquivo de DB de um teste anterior seja removido ANTES de começar.
	// Isso garante um estado 100% limpo.
	os.Remove(config.DB_URL)

	// 4. Agora, inicializa uma nova conexão de banco de dados limpa.
	db = database.InitDB()

	// 5. Cria o router e injeta a conexão do DB.
	router = gin.Default()
	routes.SetupAPI(router, db)

	// 6. Roda os testes.
	exitCode := m.Run()

	// 7. Limpeza final após todos os testes.
	os.Remove(config.DB_URL)
	os.Exit(exitCode)
}

// clearDB limpa todas as tabelas para garantir que os testes sejam independentes
func clearDB() {
	// Deleta os registros e reseta a sequência de auto incremento do SQLite
	db.Exec("DELETE FROM users")
	db.Exec("UPDATE sqlite_sequence SET seq = 0 WHERE name = 'users'")
}

// TestAuthEndpoints cobre os cenários de registro e login
func TestAuthEndpoints(t *testing.T) {
	clearDB() // Garante que o DB está limpo antes deste grupo de testes

	// --- Cenário 1: Registrar um novo usuário com sucesso ---
	t.Run("POST /auth/register - Success", func(t *testing.T) {
		registerPayload := []byte(fmt.Sprintf(`{"name": "Ada Lovelace", "email": "ada%d@example.com", "password": "password123"}`, time.Now().UnixNano()))

		req, _ := http.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(registerPayload))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Usuário registrado com sucesso", response["mensagem"])
	})

	// --- Cenário 2: Tentar registrar um usuário com e-mail já existente ---
	t.Run("POST /auth/register - Email Already Exists", func(t *testing.T) {
		uniqueEmail := fmt.Sprintf("test%d@exists.com", time.Now().UnixNano())
		userPayload := []byte(fmt.Sprintf(`{"name": "Test User", "email": "%s", "password": "password123"}`, uniqueEmail))

		// Cria o usuário uma vez
		req, _ := http.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(userPayload))
		router.ServeHTTP(httptest.NewRecorder(), req)

		// Tenta registrar com o mesmo e-mail novamente
		req, _ = http.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(userPayload))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// CORREÇÃO: Espera 409 Conflict, não 500
		assert.Equal(t, http.StatusConflict, w.Code)

		var errorResponse map[string]string
		json.Unmarshal(w.Body.Bytes(), &errorResponse)
		// CORREÇÃO: Verifica a mensagem de erro correta
		assert.Equal(t, "o e-mail informado já está em uso", errorResponse["error"])
	})

	// --- Cenário 3: Tentar registrar com dados inválidos (e-mail inválido ou senha curta) ---
	t.Run("POST /auth/register - Bad Request", func(t *testing.T) {
		registerPayload := []byte(`{"name": "Invalid User", "email": "invalid-email", "password": "short"}`)

		req, _ := http.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(registerPayload))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var errorResponse map[string]string
		json.Unmarshal(w.Body.Bytes(), &errorResponse)
		assert.Equal(t, "Dados da requisição inválidos", errorResponse["error"])
	})

	// --- Cenário 4: Fazer login com sucesso ---
	t.Run("POST /auth/login - Success", func(t *testing.T) {
		email := fmt.Sprintf("login%d@test.com", time.Now().UnixNano())
		password := "password123"
		registerPayload := []byte(fmt.Sprintf(`{"name": "Login User", "email": "%s", "password": "%s"}`, email, password))
		req, _ := http.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(registerPayload))
		router.ServeHTTP(httptest.NewRecorder(), req)

		loginPayload := []byte(fmt.Sprintf(`{"email": "%s", "password": "%s"}`, email, password))
		req, _ = http.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(loginPayload))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		// Verifica se o cookie do token foi enviado
		assert.NotEmpty(t, w.Header().Get("Set-Cookie"))
	})

	// --- Cenário 5: Tentar fazer login com credenciais inválidas ---
	t.Run("POST /auth/login - Invalid Credentials", func(t *testing.T) {
		loginPayload := []byte(`{"email": "nonexistent@user.com", "password": "wrongpassword"}`)

		req, _ := http.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(loginPayload))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		var errorResponse map[string]string
		json.Unmarshal(w.Body.Bytes(), &errorResponse)
		assert.Equal(t, "e-mail ou senha inválidos", errorResponse["error"])
	})
}

// setupAuthenticatedUser é uma função helper para criar um usuário e retornar seu ID e token
func setupAuthenticatedUser(t *testing.T) (uint, string) {
	email := fmt.Sprintf("test.user%d@example.com", time.Now().UnixNano())
	password := "password123"
	registerPayload := []byte(fmt.Sprintf(`{"name": "Test User", "email": "%s", "password": "%s"}`, email, password))
	req, _ := http.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(registerPayload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(httptest.NewRecorder(), req)

	loginPayload := []byte(fmt.Sprintf(`{"email": "%s", "password": "%s"}`, email, password))
	wLogin := httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(loginPayload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(wLogin, req)

	var authToken string
	cookies := wLogin.Result().Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "token" {
			authToken = cookie.Value
		}
	}
	assert.NotEmpty(t, authToken)

	var createdUser user.User
	db.Where("email = ?", email).First(&createdUser)
	assert.NotZero(t, createdUser.ID)

	return createdUser.ID, authToken
}

// TestUserEndpoints cobre os cenários de CRUD de usuário
func TestUserEndpoints(t *testing.T) {
	clearDB() // Garante que o DB está limpo antes deste grupo de testes
	userID, authToken := setupAuthenticatedUser(t)

	// --- Cenário 1: Listar todos os usuários com autenticação ---
	t.Run("GET /users - Success with Auth", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/users/", nil)
		req.Header.Set("Authorization", "Bearer "+authToken)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var users []user.User
		json.Unmarshal(w.Body.Bytes(), &users)
		assert.NotEmpty(t, users)
	})

	// --- Cenário 2: Tentar listar usuários sem autenticação ---
	t.Run("GET /users - Unauthorized", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/users/", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	// --- Cenário 3: Obter um usuário por ID com sucesso ---
	t.Run("GET /users/:id - Success", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/users/%d", userID), nil)
		req.Header.Set("Authorization", "Bearer "+authToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	// --- Cenário 4: Atualizar um usuário com sucesso ---
	t.Run("PUT /users/:id - Success", func(t *testing.T) {
		newName := "Updated Name"
		updatePayload := []byte(fmt.Sprintf(`{"name": "%s"}`, newName))
		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/users/%d", userID), bytes.NewBuffer(updatePayload))
		req.Header.Set("Authorization", "Bearer "+authToken)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var updatedUser user.User
		json.Unmarshal(w.Body.Bytes(), &updatedUser)
		assert.Equal(t, newName, updatedUser.Name)
	})

	// --- Cenário 5: Tentar atualizar um e-mail para um já existente ---
	t.Run("PUT /users/:id - Email Conflict", func(t *testing.T) {
		// Cria um segundo usuário para causar o conflito
		secondUserID, _ := setupAuthenticatedUser(t)
		var secondUser user.User
		db.First(&secondUser, secondUserID)

		// Tenta atualizar o primeiro usuário com o e-mail do segundo
		updatePayload := []byte(fmt.Sprintf(`{"email": "%s"}`, secondUser.Email))
		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/users/%d", userID), bytes.NewBuffer(updatePayload))
		req.Header.Set("Authorization", "Bearer "+authToken)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
	})

	// --- Cenário 6: Deletar um usuário com sucesso ---
	t.Run("DELETE /users/:id - Success", func(t *testing.T) {
		userToDeleteID, token := setupAuthenticatedUser(t) // Cria usuário novo para deletar
		url := fmt.Sprintf("/api/users/%d", userToDeleteID)

		req, _ := http.NewRequest(http.MethodDelete, url, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNoContent, w.Code)

		// Verifica se o usuário foi realmente deletado
		req, _ = http.NewRequest(http.MethodGet, url, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
