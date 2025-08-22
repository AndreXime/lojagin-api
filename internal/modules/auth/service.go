package auth

import (
	"LojaGin/internal/config"
	"LojaGin/internal/modules/user"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

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
		"sub": strconv.FormatUint(uint64(createdUser.ID), 10), // Converte ID para string
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(config.JWT_SECRET))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *Service) Login(req user.LoginUserRequest) (string, error) {
	foundUser, err := s.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(req.Password))
	if err != nil {
		return "", ErrInvalidCredentials
	}

	// Gera o token JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": strconv.FormatUint(uint64(foundUser.ID), 10), // Converte ID para string
		"exp": time.Now().Add(time.Hour * 24).Unix(),        // Token expira em 24 horas
	})

	tokenString, err := token.SignedString([]byte(config.JWT_SECRET))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
