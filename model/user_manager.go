package model

type UserManager interface {
	RegisterUser(req RegisterRequest) (*User, error)
	LoginUser(req LoginRequest) (*User, error)
}

var _ UserManager = (*UserBank)(nil)
