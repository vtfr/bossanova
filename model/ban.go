package model

import "time"

// Ban is a user's ban
type Ban struct {
	ID string `json:"id" bson:"_id"`

	IP        string        `json:"ip" bson:"ip"`
	CreatedAt time.Time     `json:"createdAt" bson:"createdAt"`
	Duration  time.Duration `json:"duration" bson:"duration"`
	Reason    string        `json:"reason" bson:"reason"`
	Mod       string        `json:"mod" bson:"mod"`
}
