package st_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vtfr/bossanova/common"
	"github.com/vtfr/bossanova/model"
)

var _ = Describe("Post", func() {

	sample := model.NewPost("", "b", "name", "subject", "comment", "ip")

	It("should create a new resource", func() {
		err := st.CreatePost(sample)

		Expect(err).To(BeNil())
	})
	It("should fail to insert a duplicated post", func() {
		err := st.CreatePost(sample)

		Expect(err).To(Equal(common.ErrConflict))
	})
	It("should retrieve a especific resource", func() {
		post, err := st.GetPost(sample.ID)

		Expect(err).To(BeNil())
		Expect(post).To(Equal(sample))
	})
	It("should be able to update a resource", func() {
		By("updating the resource")
		sample.Comment = "changed comment"
		err := st.UpdatePost(sample)

		Expect(err).To(BeNil())

		By("checking if it was updated")
		post, err := st.GetPost(sample.ID)
		Expect(err).To(BeNil())
		Expect(post).To(Equal(sample))
	})
	It("should delete a resource", func() {
		By("deleting the resource")
		err := st.DeletePost(sample.ID)

		Expect(err).To(BeNil())

		By("failling to fetch it")
		_, err = st.GetPost(sample.ID)
		Expect(err).To(Equal(common.ErrNotFound))
	})
})
