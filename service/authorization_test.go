package service_test

import (
	"testing"

	"github.com/franela/goblin"
	"github.com/vtfr/bossanova/model"
	"github.com/vtfr/bossanova/service"
)

func TestAuthorizator(t *testing.T) {
	auth := service.NewAuthorizator(map[string][]string{
		"admin":     {"a", "b"},
		"mod":       {"c", "d"},
		"anonymous": {"c"},
	})

	g := goblin.Goblin(t)
	g.Describe("Authorizator", func() {
		g.It("Should authorize correctly", func() {
			admin := &model.User{Role: "admin"}
			mod := &model.User{Role: "mod"}

			g.Assert(auth.IsAuthorized(admin, "a")).IsTrue()
			g.Assert(auth.IsAuthorized(admin, "b")).IsTrue()
			g.Assert(auth.IsAuthorized(mod, "c")).IsTrue()
			g.Assert(auth.IsAuthorized(mod, "d")).IsTrue()

			g.Assert(auth.IsAuthorized(admin, "c")).IsFalse()
			g.Assert(auth.IsAuthorized(admin, "d")).IsFalse()
			g.Assert(auth.IsAuthorized(mod, "a")).IsFalse()
			g.Assert(auth.IsAuthorized(mod, "b")).IsFalse()
		})
		g.It("Should authorize anonymous correctly", func() {
			g.Assert(auth.IsAuthorized(nil, "a")).IsFalse()
			g.Assert(auth.IsAuthorized(nil, "b")).IsFalse()
			g.Assert(auth.IsAuthorized(nil, "c")).IsTrue()
			g.Assert(auth.IsAuthorized(nil, "d")).IsFalse()
		})
	})
}
