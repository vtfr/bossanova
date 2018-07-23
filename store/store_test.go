package store_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vtfr/bossanova/common"
	"github.com/vtfr/bossanova/model"
	"github.com/vtfr/bossanova/store"
)

func TestStore(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "MongoStore")
}

var st *store.MongoStore

// Attempt to connect to store
var _ = BeforeSuite(func() {
	// if testing.Short() {
	// 	Skip("skipping Store")
	// }

	sti, err := store.NewStore("localhost:27017", "tests")
	Expect(err).To(BeNil())
	Expect(sti).NotTo(BeNil())

	// converts Store to MongoStore
	var ok bool
	st, ok = sti.(*store.MongoStore)
	Expect(ok).To(BeTrue())

	// drops the testing database so we have a clean database
	Expect(st.Database().DropDatabase()).To(BeNil())
})

var _ = AfterSuite(func() {
	st.Close()
})

var _ = Describe("MongoStore", func() {

	Context("boards", func() {
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

	Context("posts", func() {
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

	Context("users", func() {
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
})
