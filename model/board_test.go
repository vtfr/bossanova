package model_test

import (
	"testing"

	"github.com/franela/goblin"
	"github.com/vtfr/bossanova/model"
)

func TestBoard(t *testing.T) {
	g := goblin.Goblin(t)
	g.Describe("Board", func() {
		g.It("Should be valid", func() {
			boards := []struct {
				Board *model.Board
				Valid bool
			}{
				{
					Board: &model.Board{
						URI:  "b",
						Name: "Random",
					},
					Valid: true,
				},
				{
					Board: &model.Board{
						URI: "b",
					},
					Valid: false,
				},
				{
					Board: &model.Board{
						URI: "",
					},
					Valid: false,
				},
				{
					Board: &model.Board{
						URI: "0",
					},
					Valid: false,
				},
			}

			for _, b := range boards {
				if b.Valid {
					g.Assert(model.Validate(b.Board) == nil).IsTrue()
				} else {
					g.Assert(model.Validate(b.Board) == nil).IsFalse()
				}
			}
		})
	})
}
