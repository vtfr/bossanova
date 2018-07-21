package service_test

import (
	"strings"
	"testing"

	"github.com/franela/goblin"
	"github.com/vtfr/bossanova/model"
	"github.com/vtfr/bossanova/service"
)

const samplePermissionData = `{
	"roles": {
		"roleA": [
			"permissionA",
			"permissionB"
		],
		"roleB": [
			"permissionC",
			"permissionD"
		],
		"default": [
			"permissionA"
		]
	},
	"inheritance": {
		"roleA": [ "roleB" ]
	}
}`

func TestAuthorizator(t *testing.T) {
	var auth service.Authorizator

	g := goblin.Goblin(t)
	g.Describe("Authorizator", func() {
		g.It("Should warn on malformed configuration file", func() {
			r := strings.NewReader(`{
				"roles": {
					"roleA": true,
				}
			}`)

			_, err := service.NewAuthorizatorFromFile(r)
			g.Assert(err != nil).IsTrue()
		})
		g.It("Should read from the configuration file", func() {
			r := strings.NewReader(samplePermissionData)

			var err error
			auth, err = service.NewAuthorizatorFromFile(r)
			g.Assert(err == nil).IsTrue()
		})
		g.It("Should authorize correctly", func() {
			userA := &model.User{Role: "roleA"}
			userB := &model.User{Role: "roleB"}

			g.Assert(auth.IsAuthorized(userA, "permissionA")).IsTrue()
			g.Assert(auth.IsAuthorized(userA, "permissionB")).IsTrue()
			g.Assert(auth.IsAuthorized(userA, "permissionC")).IsTrue()
			g.Assert(auth.IsAuthorized(userA, "permissionD")).IsTrue()

			g.Assert(auth.IsAuthorized(userB, "permissionA")).IsFalse()
			g.Assert(auth.IsAuthorized(userB, "permissionB")).IsFalse()
			g.Assert(auth.IsAuthorized(userB, "permissionC")).IsTrue()
			g.Assert(auth.IsAuthorized(userB, "permissionD")).IsTrue()
		})
		g.It("Should authorize default correctly", func() {
			g.Assert(auth.IsAuthorized(nil, "permissionA")).IsTrue()
			g.Assert(auth.IsAuthorized(nil, "permissionB")).IsFalse()
			g.Assert(auth.IsAuthorized(nil, "permissionC")).IsFalse()
			g.Assert(auth.IsAuthorized(nil, "permissionD")).IsFalse()
		})
	})
}
