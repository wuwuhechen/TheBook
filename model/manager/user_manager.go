package manager

import "TheBook/model/structs"

type UserManager interface {
	RegisterUser(req structs.RegisterRequest) (*structs.User, error)
	LoginUser(req structs.LoginRequest) (*structs.User, error)
}
