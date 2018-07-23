package model_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vtfr/bossanova/model"
)

var _ = Describe("Model", func() {
	Context("user", func() {
		It("should create a new user", func() {
			user := model.NewUser("username", "password", "role")

			Expect(user.Username).To(Equal("username"))
			Expect(user.Valid()).To(BeNil())
			Expect(user.Role).To(Equal("role"))
		})
		It("should be able to set a password", func() {
			user := model.NewUser("username", "old", "role")
			oldHash := user.HashedPassword

			user.SetPassword("new")
			newHash := user.HashedPassword
			Expect(oldHash).NotTo(Equal(newHash))
		})
		It("should verify a password correctly", func() {
			user := model.NewUser("username", "right", "role")

			Expect(user.VerifyPassword("wrong")).To(BeFalse())
			Expect(user.VerifyPassword("right")).To(BeTrue())
		})
	})
})
