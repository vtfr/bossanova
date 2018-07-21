package service

import (
	"encoding/json"
	"io"

	"github.com/mikespook/gorbac"
	"github.com/vtfr/bossanova/model"
)

// Default role
const Default = "default"

// Authorizator determines if a user can perform an action or not
type Authorizator interface {
	IsAuthorized(user *model.User, action string) bool
}

// authorizator is the default Authorizator implementation
type authorizator struct {
	rbac *gorbac.RBAC
}

// NewAuthorizatorFromFile creates a new authorizator reading content from an
// existing permission file
func NewAuthorizatorFromFile(f io.Reader) (Authorizator, error) {
	var data struct {
		Roles       map[string][]string `json:"roles"`
		Inheritance map[string][]string `json:"inheritance"`
	}

	// decode JSON content
	if err := json.NewDecoder(f).Decode(&data); err != nil {
		return nil, err
	}

	return NewAuthorizator(data.Roles, data.Inheritance), nil
}

// NewAuthorizator creates a new Authorizator
func NewAuthorizator(roles map[string][]string,
	inheritance map[string][]string) Authorizator {

	rbac := gorbac.New()

	// create all roles and set their permissions
	for roleName := range roles {
		role := gorbac.NewStdRole(roleName)
		for _, permissionName := range roles[roleName] {
			premission := gorbac.NewStdPermission(permissionName)
			role.Assign(premission)
		}

		rbac.Add(role)
	}

	// set inheritances
	for role := range inheritance {
		rbac.SetParents(role, inheritance[role])
	}

	return &authorizator{rbac}
}

// IsAuthorized returns true if a user can perform a certain action, else
// return false
func (auth *authorizator) IsAuthorized(user *model.User, action string) bool {
	role := Default
	if user != nil {
		role = user.Role
	}

	return auth.rbac.IsGranted(role, gorbac.NewStdPermission(action), nil)
}
