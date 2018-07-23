package st_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vtfr/bossanova/common"
	"github.com/vtfr/bossanova/model"
)

var _ = Describe("Board", func() {

	sample := model.NewBoard("b", "name", "description")

	It("should create a new resource", func() {
		err := st.CreateBoard(sample)

		Expect(err).To(BeNil())
	})
	It("should fail to insert a duplicated board", func() {
		err := st.CreateBoard(sample)

		Expect(err).To(Equal(common.ErrConflict))
	})
	It("should retrieve all the resources", func() {
		boards, err := st.AllBoards()

		Expect(err).To(BeNil())
		Expect(boards).To(HaveLen(1))
		Expect(boards).To(ContainElement(sample))
	})
	It("should retrieve a especific resource", func() {
		board, err := st.GetBoard(sample.URI)

		Expect(err).To(BeNil())
		Expect(board).To(Equal(sample))
	})
	It("should be able to update a resource", func() {
		By("updating the resource")
		sample.Description = "changed description"
		err := st.UpdateBoard(sample)

		Expect(err).To(BeNil())

		By("checking if it was updated")
		board, err := st.GetBoard(sample.URI)
		Expect(err).To(BeNil())
		Expect(board).To(Equal(sample))
	})
	It("should delete a resource", func() {
		By("deleting the resource")
		err := st.DeleteBoard(sample.URI)

		Expect(err).To(BeNil())

		By("failling to fetch it")
		_, err = st.GetBoard(sample.URI)
		Expect(err).To(Equal(common.ErrNotFound))
	})
})
