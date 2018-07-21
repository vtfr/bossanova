package model_test

import (
	"github.com/franela/goblin"
	"github.com/vtfr/bossanova/model"
	"testing"
)

func TestUser(t *testing.T) {
	g := goblin.Goblin(t)
	g.Describe("User", func() {
		g.It("Should set and verify a password sucessfully", func() {
			user := &model.User{}
			user.SetPassword("right")

			g.Assert(user.VerifyPassword("wrong")).IsFalse("wrong password")
			g.Assert(user.VerifyPassword("right")).IsTrue("right password")
		})
	})
}
