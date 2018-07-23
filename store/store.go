// Package store contains all permanent store related tasks, such as creating,
// retrieving, updating and deleting models from the underlying database.
package store

import (
	"github.com/globalsign/mgo"
)

// Store abstracts all database access
//go:generate mockgen -destination=../mocks/store.go -package=mocks github.com/vtfr/bossanova/store Store
type Store interface {
	BoardStore
	PostStore
	UserStore
	BanStore

	Clone() Store
	Close()
}

// MongoStore is a MongoDB implementation of Store
type MongoStore struct {
	session  *mgo.Session
	database *mgo.Database
}

// NewStore create a new MongoDB backed data Store
func NewStore(uri string, database string) (Store, error) {
	ses, err := mgo.Dial(uri)
	if err != nil {
		return nil, err
	}

	return &MongoStore{
		session:  ses,
		database: ses.DB(database),
	}, nil
}

// Session returns the internal MongoDB session
func (s *MongoStore) Session() *mgo.Session {
	return s.session
}

// Database returns the internal MongoDB database
func (s *MongoStore) Database() *mgo.Database {
	return s.database
}

// Boards return the boards collection
func (s *MongoStore) Boards() *mgo.Collection {
	return s.database.C("boards")
}

// Posts return the posts collection
func (s *MongoStore) Posts() *mgo.Collection {
	return s.database.C("posts")
}

// Users return the users collection
func (s *MongoStore) Users() *mgo.Collection {
	return s.database.C("users")
}

// Bans return the bans collection
func (s *MongoStore) Bans() *mgo.Collection {
	return s.database.C("bans")
}

// Clone creates a new Store with a new MongoDB connection
func (s *MongoStore) Clone() Store {
	ses := s.session.Copy()
	return &MongoStore{
		session:  ses,
		database: ses.DB(s.database.Name),
	}
}

// Close closes a MongoDB connection
func (s *MongoStore) Close() {
	s.session.Close()
}
