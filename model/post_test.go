package model_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vtfr/bossanova/model"
)

var _ = Describe("Post", func() {
	It("should create a new thread sucessfully", func() {
		post := model.NewPost("", "board", "name", "subject", "comment", "ip")

		Expect(post.IsReply()).To(BeFalse())
		Expect(post.Board).To(Equal("board"))
		Expect(post.Name).To(Equal("name"))
		Expect(post.Subject).To(Equal("subject"))
		Expect(post.Comment).To(Equal("comment"))
		Expect(post.IP).To(Equal("ip"))
		Expect(post.CreatedAt.IsZero()).To(BeFalse())
		Expect(post.LastBumpedAt).ToNot(BeNil())
		Expect(post.Valid()).ToNot(HaveOccurred())
	})
	It("should create a new reply sucessfully", func() {
		post := model.NewPost("parent", "board", "name", "subject", "comment", "ip")

		Expect(post.IsReply()).To(BeTrue())
		Expect(post.Parent).To(Equal("parent"))
		Expect(post.Board).To(Equal("board"))
		Expect(post.Name).To(Equal("name"))
		Expect(post.Subject).To(Equal("subject"))
		Expect(post.Comment).To(Equal("comment"))
		Expect(post.IP).To(Equal("ip"))
		Expect(post.CreatedAt.IsZero()).To(BeFalse())
		Expect(post.LastBumpedAt).To(BeNil())
		Expect(post.Valid()).ToNot(HaveOccurred())
	})

	// TODO implement business logic checking
})
