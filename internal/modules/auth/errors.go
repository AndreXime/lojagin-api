package auth

import "errors"

var (
	ErrInvalidCredentials = errors.New("e-mail ou senha inválidos")
)
