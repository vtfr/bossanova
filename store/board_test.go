package store_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vtfr/bossanova/common"
	"github.com/vtfr/bossanova/model"
)

var _ = Describe("Board", func() {

	sample := model.NewBoard("b", "name", "description")

	It("should create a new resource", func() {
		SkipIfShort()

		err := st.CreateBoard(sample)

		Expect(err).NotTo(HaveOccurred())
	})
	It("should fail to insert a duplicated board", func() {
		SkipIfShort()

		err := st.CreateBoard(sample)

		Expect(err).To(Equal(common.ErrConflict))
	})
	It("should retrieve all the resources", func() {
		SkipIfShort()

		boards, err := st.AllBoards()

		Expect(err).NotTo(HaveOccurred())
		Expect(boards).To(HaveLen(1))
		Expect(boards).To(ContainElement(sample))
	})
	It("should retrieve a especific resource", func() {
		SkipIfShort()

		board, err := st.GetBoard(sample.URI)

		Expect(err).NotTo(HaveOccurred())
		Expect(board).To(Equal(sample))
	})
	It("should be able to update a resource", func() {
		SkipIfShort()

		By("updating the resource")
		sample.Description = "changed description"
		err := st.UpdateBoard(sample)

		Expect(err).NotTo(HaveOccurred())

		By("checking if it was updated")
		board, err := st.GetBoard(sample.URI)
		Expect(err).NotTo(HaveOccurred())
		Expect(board).To(Equal(sample))
	})
	It("should delete a resource", func() {
		SkipIfShort()

		By("deleting the resource")
		err := st.DeleteBoard(sample.URI)

		Expect(err).NotTo(HaveOccurred())

		By("failling to fetch it")
		_, err = st.GetBoard(sample.URI)
		Expect(err).To(Equal(common.ErrNotFound))
	})
})
