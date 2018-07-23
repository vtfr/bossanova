package model_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vtfr/bossanova/model"
)

var _ = Describe("Model", func() {
	Context("board", func() {
		It("should create a new board sucessfully", func() {
			board := model.NewBoard("uri", "name", "description")

			Expect(board.URI).To(Equal("uri"))
			Expect(board.Name).To(Equal("name"))
			Expect(board.Description).To(Equal("description"))
			Expect(board.CreatedAt.IsZero()).To(BeFalse())
			Expect(board.Valid()).To(BeNil())
		})
		It("should have a valid URI", func() {
			By("creating a invalid boards")
			Expect(model.NewBoard("", "name", "description").Valid()).ToNot(BeNil())
			Expect(model.NewBoard("1", "name", "description").Valid()).ToNot(BeNil())
			Expect(model.NewBoard("^d", "name", "description").Valid()).ToNot(BeNil())

			By("creating a valid board")
			Expect(model.NewBoard("uri", "name", "description").Valid()).To(BeNil())
		})
		It("should have a name", func() {
			By("creating a invalid board")
			Expect(model.NewBoard("uri", "", "description").Valid()).ToNot(BeNil())

			By("creating a valid board")
			Expect(model.NewBoard("uri", "name", "description").Valid()).To(BeNil())
		})
	})

})
