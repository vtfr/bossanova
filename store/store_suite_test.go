package store_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vtfr/bossanova/store"
)

var st *store.MongoStore

func TestStore(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Store Suite")
}

var _ = BeforeSuite(func() {
	if !testing.Short() {
		sti, err := store.NewStore("mongodb://localhost:27017", "tests")
		if err != nil {
			Fail("Failed connecting to local database. Skipping tests")
		}

		st = sti.(*store.MongoStore)
		st.Database().DropDatabase()
	}
})

// WhenTestingStoreIt is used to skip tests if no database is set
func SkipIfShort() {
	if testing.Short() {
		Skip("Skipping because testing is short")
	}
}

var _ = AfterSuite(func() {
	if !testing.Short() {
		st.Database().DropDatabase()
		st.Close()
	}
})
