package st_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vtfr/bossanova/common"
	"github.com/vtfr/bossanova/model"
)

var _ = Describe("User", func() {

	sample := model.NewUser("username", "password", "admin")

	It("should create a new resource", func() {
		err := st.CreateUser(sample)

		Expect(err).To(BeNil())
	})
	It("should fail to insert a duplicated resource", func() {
		err := st.CreateUser(sample)

		Expect(err).To(Equal(common.ErrConflict))
	})
	It("should retrieve all the resources", func() {
		users, err := st.AllUsers()

		Expect(err).To(BeNil())
		Expect(users).To(HaveLen(1))
		Expect(users).To(ContainElement(sample))
	})
	It("should retrieve a especific resource", func() {
		user, err := st.GetUser(sample.Username)

		Expect(err).To(BeNil())
		Expect(user).To(Equal(sample))
	})
	It("should be able to update a resource", func() {
		By("updating the resource")
		sample.Role = "mod"
		err := st.UpdateUser(sample)

		Expect(err).To(BeNil())

		By("checking if it was updated")
		user, err := st.GetUser(sample.Username)
		Expect(err).To(BeNil())
		Expect(user).To(Equal(sample))
	})
	It("should delete a resource", func() {
		By("deleting the resource")
		err := st.DeleteUser(sample.Username)

		Expect(err).To(BeNil())

		By("failling to fetch it")
		_, err = st.GetUser(sample.Username)
		Expect(err).To(Equal(common.ErrNotFound))
	})
})
