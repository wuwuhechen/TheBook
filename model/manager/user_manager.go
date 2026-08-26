package manager

import "TheBook/model/structs"

type UserManager interface {
	RegisterUser(req structs.RegisterRequest) (*structs.User, error)
	LoginUser(req structs.LoginRequest) (*structs.User, error)
	GetUser(username string) (*structs.User, bool)
	// SetPasswordHash(username, passwordHash string)
	AddUser(user *structs.User) error
	DeleteUser(username string) error
}
