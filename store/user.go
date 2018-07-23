package store

import (
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/vtfr/bossanova/model"
)

// UserStore abstracts user data access
type UserStore interface {
	// CreateUser creates a user
	CreateUser(user *model.User) error
	// AllUsers return all users
	AllUsers() ([]*model.User, error)
	// GetUset gets an user
	GetUser(username string) (*model.User, error)
	// UpdateUser updates a user
	UpdateUser(user *model.User) error
	// DeleteUser deletes a user
	DeleteUser(username string) error
}

// CreateUser creates a user
func (s *MongoStore) CreateUser(user *model.User) error {
	if user.CreatedAt.IsZero() {
		user.CreatedAt = time.Now()
	}

	if user.LastModifiedAt.IsZero() {
		user.LastModifiedAt = user.CreatedAt
	}

	return mgoErr(s.Users().Insert(user))
}

// AllUsers return all users
func (s *MongoStore) AllUsers() ([]*model.User, error) {
	users := []*model.User{}
	return users, mgoErr(s.Users().Find(nil).All(&users))
}

// GetUset gets an user
func (s *MongoStore) GetUser(username string) (*model.User, error) {
	user := &model.User{}
	return user, mgoErr(s.Users().FindId(username).One(&user))
}

// UpdateUser updates a user
func (s *MongoStore) UpdateUser(user *model.User) error {
	return mgoErr(s.Users().UpdateId(user.Username, bson.M{
		"$set": user,
	}))
}

// DeleteUser deletes a user
func (s *MongoStore) DeleteUser(username string) error {
	return mgoErr(s.Users().RemoveId(username))
}
