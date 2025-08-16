package user

import (
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo *UserRepository
}

func NewService(r *UserRepository) *Service {
	return &Service{repo: r}
}

func (s *Service) FindByID(id int64) (*User, error) {
	return s.repo.GetUserByID(id)
}

func (s *Service) FindAll() ([]User, error) {
	return s.repo.GetAllUsers()
}

func (s *Service) Update(id int64, req UpdateUserRequest) (*User, error) {
	userToUpdate, err := s.repo.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		userToUpdate.Name = *req.Name
	}
	if req.Email != nil {
		if *req.Email != userToUpdate.Email {
			taken, err := s.repo.IsEmailTaken(*req.Email, id)
			if err != nil {
				return nil, err
			}
			if taken {
				return nil, ErrEmailExists
			}
			userToUpdate.Email = *req.Email
		}
	}
	if req.Password != nil {
		pw := strings.TrimSpace(*req.Password)
		switch {
		case pw == "":
			return nil, ErrInvalidPassword
		case len(pw) < 8:
			return nil, ErrShortPassword
		case len(pw) > 72:
			return nil, ErrLongPassword
		}
		hashed, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		userToUpdate.Password = string(hashed)
	}

	if err := s.repo.UpdateUser(userToUpdate); err != nil {
		return nil, err
	}
	return userToUpdate, nil
}

func (s *Service) Delete(id int64) error {
	return s.repo.DeleteUser(id)
}
