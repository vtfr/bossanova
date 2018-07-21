package store

import "github.com/vtfr/bossanova/model"

// UserStore abstracts user data access
type UserStore interface {
	AllUsers() ([]*model.User, error)
	GetUser(username string) (*model.User, error)
	CreateUser(user *model.User) error
	UpdateUser(user *model.User) error
	DeleteUser(username string) error
}
