package service

import (
	"io"
	"net/http"

	"github.com/vtfr/bossanova/common"
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
	ErrBanned         = common.NewDetailedError("banned", "you are banned", http.StatusForbidden)
	ErrThreadNotFound = common.NewDetailedError("not-found", "thread not found", http.StatusNotFound)
	ErrBoardNotFound  = common.NewDetailedError("not-found", "board not found", http.StatusNotFound)
	ErrThreadLocked   = common.NewDetailedError("locked", "thread locked", http.StatusForbidden)
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
	if err != nil && err != store.ErrNotFound {
		return nil, err
	}

	// if banned and can't ignore ban, then return error
	if banned && !auth.IsAuthorized(user, "ban.ignore") {
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
		if thread.Locked && !auth.IsAuthorized(user, "post.ingnore-locked") {
			return nil, ErrThreadLocked
		}
	}

	// create post
	return post, s.CreatePost(post)
}
