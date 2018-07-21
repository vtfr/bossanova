package model

import (
	"crypto/md5"
	"encoding/base64"
	"strconv"
	"time"
)

// Post is a comment
type Post struct {
	ID     string `json:"id" bson:"_id"`
	Parent string `json:"parent,omitempty" bson:"parent,omitempty"`
	Board  string `json:"board,omitempty" bson:"board,omitempty" validate:"board,required"`

	Name    string `json:"name" bson:"name" validate:"max=16"`
	Comment string `json:"comment" bson:"comment" validate:"required,min=4,max=5120"`
	Subject string `json:"subject,omitempty" bson:"subject,omitempty" validate:"max=48"`
	Deleted bool   `json:"deleted,omitempty" bson:"deleted"`
	IP      string `json:"address,omitempty" bson:"address"`

	Locked bool `json:"locked,omitempty" bson:"locked"`
	Hidden bool `json:"hidden,omitempty" hidden:"locked"`

	Medias []*Media `json:"medias" bson:"medias"`

	CreatedAt    time.Time  `json:"createdAt" bson:"createdAt"`
	LastBumpedAt *time.Time `json:"lastBumpedAt,omitempty" bson:"lastBumpedAt"`
}

// NewPost creates a new post
func NewPost(parent, board, name, subject, comment, ip string) *Post {
	now := time.Now()
	id := GeneratePostID(now, comment, subject)

	post := &Post{
		ID:        id,
		Parent:    parent,
		Board:     board,
		Name:      name,
		Subject:   subject,
		Comment:   comment,
		Medias:    []*Media{},
		IP:        ip,
		CreatedAt: now,
	}

	// if OP, add LastBumpedAt field
	if !post.IsReply() {
		post.LastBumpedAt = &now
	}

	return post
}

// AddMedia adds a new media
func (p *Post) AddMedia(m *Media) {
	p.Medias = append(p.Medias, m)
}

// IsReply returns whether a post is a reply or not
func (p *Post) IsReply() bool {
	return p.Parent != ""
}

// Thread is a structure holding a collection of posts, where one of them is
// the common parent (OP) and the others are replies
type Thread struct {
	Post    `json:",inline" bson:",inline"`
	Replies []*Post `json:"replies" bson:"replies"`
}

// GeneratePostID creates a new post ID based on the post's subject, comment and
// creation time
func GeneratePostID(createdAt time.Time, comment, subject string) string {
	hasher := md5.New()
	hasher.Write([]byte(strconv.FormatInt(createdAt.Unix(), 10)))
	hasher.Write([]byte(comment))
	hasher.Write([]byte(subject))

	// Generate the base64 encoded hash truncated at 6 runes
	hash := base64.RawURLEncoding.EncodeToString(hasher.Sum(nil))
	return hash[:6]
}
