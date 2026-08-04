package service

import "TheBook/model"

type UserServer interface {
	RegisterUser(req model.RegisterRequest) (*model.User, error)
	LoginUser(req model.LoginRequest) (*model.User, error)
}

var _ UserServer = (*model.UserBank)(nil)
