package auth

import (
	"LojaGin/internal/user"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

type Service struct {
	userRepo *user.UserRepository
}

func NewService(userRepo *user.UserRepository) *Service {
	return &Service{userRepo: userRepo}
}

func (s *Service) Register(req user.CreateUserRequest) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	newUser := &user.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	createdUser, err := s.userRepo.CreateUser(newUser)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": createdUser.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *Service) Login(req user.LoginUserRequest) (string, error) {
	foundUser, err := s.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		return "", errors.New("credenciais inválidas")
	}

	err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(req.Password))
	if err != nil {
		return "", errors.New("credenciais inválidas")
	}

	// Gera o token JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": foundUser.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(), // Token expira em 24 horas
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
