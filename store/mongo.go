package store

import (
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/sirupsen/logrus"
	"github.com/vtfr/bossanova/common"
	"github.com/vtfr/bossanova/model"
)

// MongoStore is a MongoDB implementation of Store
type MongoStore struct {
	Session  *mgo.Session
	Database *mgo.Database

	// Collections
	usersColl  *mgo.Collection
	postsColl  *mgo.Collection
	boardsColl *mgo.Collection
	bansColl   *mgo.Collection
}

// Compile-time check if MongoStore implements Store
var _ Store = &MongoStore{}

// NewMongoStore creates a new MongoStore
func NewMongoStore(uri string) (Store, error) {
	ses, err := mgo.Dial(uri)
	if err != nil {
		return nil, err
	}

	return newMongoStore(ses), nil
}

// Clone creates a new MongoDB connection
func (s *MongoStore) Clone() Store {
	return newMongoStore(s.Session.Copy())
}

// Close closes a MongoDB connection
func (s *MongoStore) Close() {
	s.Session.Close()
}

func newMongoStore(ses *mgo.Session) *MongoStore {
	db := ses.DB("bossanova")
	return &MongoStore{
		Session:  ses,
		Database: db,

		boardsColl: db.C("boards"),
		postsColl:  db.C("posts"),
		usersColl:  db.C("users"),
		bansColl:   db.C("bans"),
	}
}

// Boards

func (s *MongoStore) CreateBoard(board *model.Board) error {
	return mgoErr(s.boardsColl.Insert(&board))
}

func (s *MongoStore) AllBoards() ([]*model.Board, error) {
	boards := []*model.Board{}
	return boards, mgoErr(s.boardsColl.Find(nil).All(&boards))
}

func (s *MongoStore) GetBoard(uri string) (board *model.Board, err error) {
	return board, mgoErr(s.boardsColl.FindId(uri).One(&board))
}

func (s *MongoStore) UpdateBoard(board *model.Board) error {
	return mgoErr(s.boardsColl.UpdateId(board.URI, board))
}

func (s *MongoStore) DeleteBoard(uri string) error {
	return mgoErr(s.boardsColl.RemoveId(uri))
}

func (s *MongoStore) CreatePost(p *model.Post) error {
	// Assign local-time if no time is found
	if p.CreatedAt.IsZero() {
		p.CreatedAt = time.Now()
	}

	// Assign an ID if no ID is found
	if p.ID == "" {
		p.ID = model.GeneratePostID(p.CreatedAt, p.Comment, p.Subject)
	}

	// Insers in the store
	return mgoErr(s.postsColl.Insert(p))
}

func (s *MongoStore) GetPost(postID string) (*model.Post, error) {
	post := &model.Post{}
	return post, mgoErr(s.postsColl.FindId(postID).One(&post))
}

func (s *MongoStore) UpdatePost(p *model.Post) error {
	return mgoErr(s.postsColl.UpdateId(p.ID, p))
}

func (s *MongoStore) DeletePost(id string) error {
	return mgoErr(s.usersColl.Remove(bson.M{
		"$or": bson.M{
			"_id":    id,
			"parent": id,
		},
	}))
}

func (s *MongoStore) GetReplies(threadID string) ([]*model.Post, error) {
	replies := []*model.Post{}
	return replies, mgoErr(s.postsColl.Find(bson.M{
		"parent": threadID,
	}).All(&replies))
}

func (s *MongoStore) GetThread(threadID string) (*model.Thread, error) {
	thread := &model.Thread{}

	return thread, mgoErr(s.postsColl.Pipe([]bson.M{
		{
			"$match": bson.M{
				"_id":    threadID,
				"parent": bson.M{"$exists": false},
			},
		}, {
			"$lookup": bson.M{
				"from":         s.postsColl.Name,
				"localField":   "_id",
				"foreignField": "parent",
				"as":           "replies",
			},
		}, {
			"$sort": bson.M{
				"replies.createdAt": -1,
			},
		},
	}).One(&thread))
}

func (s *MongoStore) BumpThread(threadID string) error {
	return mgoErr(s.postsColl.UpdateId(threadID, bson.M{
		"lastBumpedAt": time.Now(),
	}))
}

func (s *MongoStore) CreateUser(user *model.User) error {
	if user.CreatedAt.IsZero() {
		user.CreatedAt = time.Now()
	}

	if user.LastModifiedAt.IsZero() {
		user.LastModifiedAt = user.CreatedAt
	}

	return mgoErr(s.usersColl.Insert(user))
}

func (s *MongoStore) AllUsers() ([]*model.User, error) {
	users := []*model.User{}
	return users, mgoErr(s.usersColl.Find(nil).All(&users))
}

func (s *MongoStore) GetUser(username string) (*model.User, error) {
	user := &model.User{}
	return user, mgoErr(s.usersColl.FindId(username).One(&user))
}

func (s *MongoStore) UpdateUser(user *model.User) error {
	return mgoErr(s.usersColl.UpdateId(user, bson.M{
		"$set": user,
	}))
}

func (s *MongoStore) DeleteUser(username string) error {
	return mgoErr(s.usersColl.RemoveId(username))
}

func (s *MongoStore) IsBanned(ip string) (ban *model.Ban, exists bool, err error) {
	err = mgoErr(s.bansColl.Find(bson.M{"ip": ip}).One(&ban))
	exists = err == nil
	return
}

// mgoErr converts mgo errors to a better error type
func mgoErr(err error) error {
	switch {
	case err == nil:
		return nil
	case err == mgo.ErrNotFound:
		return common.ErrNotFound
	case mgo.IsDup(err):
		return common.ErrConflict
	default:
		logrus.Panicln("Unknown error in MongoDB:", err)
		return nil
	}
}
