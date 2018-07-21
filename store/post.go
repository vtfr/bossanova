package store

import "github.com/vtfr/bossanova/model"

// PostStore stores all post persistent data
type PostStore interface {
	GetPost(postID string) (*model.Post, error)
	CreatePost(p *model.Post) error
	UpdatePost(p *model.Post) error
	DeletePost(postID string) error

	GetReplies(postID string) ([]*model.Post, error)
	GetThread(threadID string) (*model.Thread, error)
	BumpThread(threadID string) error
}
