package store

import (
	"github.com/globalsign/mgo/bson"
	"github.com/vtfr/bossanova/model"
)

// BanStore abstracts ban data access
type BanStore interface {
	// IsBanned verifies if an IP address is banned or not
	IsBanned(ip string) (*model.Ban, bool, error)

	// // AllValidBans return all valid Bans
	// AllValidBans() ([]*model.Ban, error)
	// // GetBan gets an specific Banww
	// GetBan(id string) (*model.Ban, error)
	// // UpdateBan updates a Ban
	// UpdateBan(Ban *model.Ban) error
	// // CreateBan creates a Ban
	// CreateBan(Ban *model.Ban) error
	// // DeleteBan deletes a Ban
	// DeleteBan(id string) error
}

// IsBanned verifies if an IP address is banned or not
func (s *MongoStore) IsBanned(ip string) (ban *model.Ban, exists bool, err error) {
	err = mgoErr(s.Bans().Find(bson.M{"ip": ip}).One(&ban))
	exists = err == nil
	return
}
