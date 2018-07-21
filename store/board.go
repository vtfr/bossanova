package store

import "github.com/vtfr/bossanova/model"

// BoardStore stores all board persistent data
type BoardStore interface {
	// AllBoards return all existing boards
	AllBoards() ([]*model.Board, error)
	// GetBoard gets an specific board
	GetBoard(uri string) (*model.Board, error)
	// UpdateBoard updates a board
	UpdateBoard(board *model.Board) error
	// CreateBoard creates a board
	CreateBoard(board *model.Board) error
	// DeleteBoard deletes a board
	DeleteBoard(uri string) error
}
