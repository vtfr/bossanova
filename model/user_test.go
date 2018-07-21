package model_test

import (
	"testing"

	"github.com/franela/goblin"
	"github.com/vtfr/bossanova/model"
)

func TestUser(t *testing.T) {
	user := &model.User{}
	user.SetPassword("right")

	g := goblin.Goblin(t)
	g.Describe("User", func() {
		g.It("Should verify a password successfully", func() {
			g.Assert(user.VerifyPassword("right")).IsTrue("right password")
		})
		g.It("Should not verify a wrong password", func() {
			g.Assert(user.VerifyPassword("wrong")).IsFalse("wrong password")
		})
	})
}
