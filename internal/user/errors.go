package user

import "errors"

var (
	ErrEmailExists     = errors.New("e-mail já cadastrado")
	ErrUserNotFound    = errors.New("usuário não encontrado")
	ErrInvalidPassword = errors.New("password cannot be empty")
	ErrShortPassword   = errors.New("password must be at least 8 characters")
	ErrLongPassword    = errors.New("password too long")
)
