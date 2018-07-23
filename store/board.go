package store

import "github.com/vtfr/bossanova/model"

// BoardStore abstracts board data access
type BoardStore interface {
	// CreateBoard creates a board
	CreateBoard(board *model.Board) error
	// AllBoards return all existing boards
	AllBoards() ([]*model.Board, error)
	// GetBoard gets an specific board
	GetBoard(uri string) (*model.Board, error)
	// UpdateBoard updates a board
	UpdateBoard(board *model.Board) error
	// DeleteBoard deletes a board
	DeleteBoard(uri string) error
}

// CreateBoard creates a board
func (s *MongoStore) CreateBoard(board *model.Board) error {
	return mgoErr(s.Boards().Insert(&board))
}

// AllBoards return all existing boards
func (s *MongoStore) AllBoards() ([]*model.Board, error) {
	boards := []*model.Board{}
	return boards, mgoErr(s.Boards().Find(nil).All(&boards))
}

// GetBoard gets an specific board
func (s *MongoStore) GetBoard(uri string) (board *model.Board, err error) {
	return board, mgoErr(s.Boards().FindId(uri).One(&board))
}

// UpdateBoard updates a board
func (s *MongoStore) UpdateBoard(board *model.Board) error {
	return mgoErr(s.Boards().UpdateId(board.URI, board))
}

// DeleteBoard deletes a board
func (s *MongoStore) DeleteBoard(uri string) error {
	return mgoErr(s.Boards().RemoveId(uri))
}
