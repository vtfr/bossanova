package service

import (
	"github.com/vtfr/bossanova/model"
)

// Anonymous role
const Anonymous = "anonymous"

// Authorizator handles all authorization tasks
type Authorizator map[string][]string

// NewAuthorizator creates a new Authorizator
func NewAuthorizator(perms map[string][]string) Authorizator {
	return Authorizator(perms)
}

// Roles return all existing roles
func (auth Authorizator) Roles() []string {
	roles := make([]string, len(auth), len(auth))

	i := 0
	for role := range auth {
		roles[i] = role
		i++
	}

	return roles
}

// IsAuthorized returns true if a user can perform a certain action, else
// return false
func (auth Authorizator) IsAuthorized(user *model.User, action string) bool {
	role := Anonymous
	if user != nil {
		role = user.Role
	}

	// search if permission exists
	for _, perm := range auth[role] {
		if perm == action {
			return true
		}
	}

	return false
}
