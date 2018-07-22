package service

import (
	"errors"
	"io"

	"github.com/vtfr/bossanova/model"
	"github.com/vtfr/bossanova/store"
)

type PostingError struct {
	message string
	status  int
}

func (p *PostingError) Error() string {
	return p.message
}

var (
	ErrBanned         = errors.New("banned")
	ErrThreadNotFound = errors.New("not-found")
	ErrBoardNotFound  = errors.New("not-found")
	ErrThreadLocked   = errors.New("locked")
)

// CreatePostRequest constains all post creation information necessary for the
// creation of a new post
type CreatePostData struct {
	Parent string `form:"parent"`
	Board  string `form:"board"`

	Name    string `form:"name"`
	Comment string `form:"comment"`
	Subject string `form:"subject"`

	User *model.User

	Medias []*CreatePostMediaData

	IP string
}

// CreatePostMediaRequest contains all media information for upload
type CreatePostMediaData struct {
	Size int64
	Name string
	File io.Reader
}

// CreatePost creates a new post
func CreatePost(auth Authorizator,
	s store.Store,
	user *model.User,
	data *CreatePostData) (*model.Post, error) {

	// creates the Post structure
	post := model.NewPost(data.Parent, data.Board, "Anonymous",
		data.Subject, data.Comment, data.IP)

	// validates post
	if err := model.Validate(post); err != nil {
		return nil, err
	}

	// verifies if user is banned
	_, banned, err := s.IsBanned(data.IP)
	if err != nil {
		return nil, err
	}

	// if banned and can't ignore ban, then return error
	if banned && !auth.IsAuthorized(user, "posting.ignore-ban") {
		return nil, ErrBanned
	}

	// verify if board exists
	board, err := s.GetBoard(post.Board)
	if err != nil {
		return nil, ErrBoardNotFound
	}

	// verify if board hidden and can post on board
	if board.Hidden && !auth.IsAuthorized(user, "board.hidden") {
		return nil, ErrBoardNotFound
	}

	if post.IsReply() {
		// verifies if thread exists
		thread, err := s.GetPost(post.Parent)
		if err != nil {
			return nil, ErrThreadNotFound
		}

		// verify if thread is actually a thread
		if thread.IsReply() {
			return nil, ErrThreadNotFound
		}

		// verify if thread locked and can bypass that lock
		if thread.Locked && !auth.IsAuthorized(user, "posting.ingnore-locked") {
			return nil, ErrThreadLocked
		}
	}

	// create post
	return post, s.CreatePost(post)
}
