package model_test

import (
	"github.com/franela/goblin"
	"github.com/vtfr/bossanova/model"
	"testing"
	"time"
)

// GenerateThreadID creates a new ThreadID based on the thread's subject,
// comment and local time
func TestGeneratePostID(t *testing.T) {
	g := goblin.Goblin(t)
	g.Describe("GeneratePostID", func() {
		g.It("Should generate the same hash for multiple calls with the same arguments", func() {
			t := time.Now()

			h1 := model.GeneratePostID(t, "comment", "subject")
			h2 := model.GeneratePostID(t, "comment", "subject")

			g.Assert(h1 == h2).IsTrue()
		})
		g.It("Should generate different hashes for different arguments", func() {
			t := time.Now()

			h1 := model.GeneratePostID(t, "comment one", "")
			h2 := model.GeneratePostID(t, "comment two", "")

			g.Assert(h1 == h2).IsFalse()
		})
	})
}
