package store

import (
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/vtfr/bossanova/model"
)

// PostStore abstracts post data access
type PostStore interface {
	// CreatePost creates a new post
	CreatePost(p *model.Post) error
	// GetPost returns a specific post
	GetPost(postID string) (*model.Post, error)
	// UpdatePost updates a post
	UpdatePost(p *model.Post) error
	// DeletePost deletes a post
	DeletePost(postID string) error

	// GetThread returns a thread with it's replies
	GetThread(threadID string) (*model.Thread, error)
	// GetReplies return all thread replies
	GetReplies(postID string) ([]*model.Post, error)
	// BumpThreads bumps a thread
	BumpThread(threadID string) error
}

// CreatePost creates a new post
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
	return mgoErr(s.Posts().Insert(p))
}

// GetPost returns a specific post
func (s *MongoStore) GetPost(postID string) (*model.Post, error) {
	post := &model.Post{}
	return post, mgoErr(s.Posts().FindId(postID).One(&post))
}

// UpdatePost updates a post
func (s *MongoStore) UpdatePost(p *model.Post) error {
	return mgoErr(s.Posts().UpdateId(p.ID, p))
}

// DeletePost deletes a post
func (s *MongoStore) DeletePost(id string) error {
	return mgoErr(s.Posts().RemoveId(id))
}

// GetThread returns a thread with it's replies
func (s *MongoStore) GetThread(threadID string) (*model.Thread, error) {
	thread := &model.Thread{}

	return thread, mgoErr(s.Posts().Pipe([]bson.M{
		{
			"$match": bson.M{
				"_id":    threadID,
				"parent": bson.M{"$exists": false},
			},
		}, {
			"$lookup": bson.M{
				"from":         s.Posts().Name,
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

// GetReplies return all thread replies
func (s *MongoStore) GetReplies(threadID string) ([]*model.Post, error) {
	replies := []*model.Post{}
	return replies, mgoErr(s.Posts().Find(bson.M{
		"parent": threadID,
	}).All(&replies))
}

// BumpThreads bumps a thread
func (s *MongoStore) BumpThread(threadID string) error {
	return mgoErr(s.Posts().UpdateId(threadID, bson.M{
		"lastBumpedAt": time.Now(),
	}))
}
