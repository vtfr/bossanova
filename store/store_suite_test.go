package st_test

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
	if testing.Short() {
		Skip("Skipping suite as tests are short")
	}

	st, err := st.NewStore("mongodb://localhost:27017", "tests")
	if err != nil {
		Skip("Failed connecting to local database. Skipping tests")
	}

	st = st.(*st.MongoStore)
	st.Database().DropDatabase()
})

var _ = AfterSuite(func() {
	st.Database().DropDatabase()
	st.Close()
})
