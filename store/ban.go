package store

import "github.com/vtfr/bossanova/model"

// BanStore stores all Ban persistent data
type BanStore interface {
	IsBanned(ip string) (*model.Ban, bool, error)

	// // AllValidBans return all valid Bans
	// AllValidBans() ([]*model.Ban, error)
	// // GetBan gets an specific Ban
	// GetBan(id string) (*model.Ban, error)
	// // UpdateBan updates a Ban
	// UpdateBan(Ban *model.Ban) error
	// // CreateBan creates a Ban
	// CreateBan(Ban *model.Ban) error
	// // DeleteBan deletes a Ban
	// DeleteBan(id string) error
}
