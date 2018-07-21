package store

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
