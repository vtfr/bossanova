package service_test

import (
	"errors"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vtfr/bossanova/mocks"
	"github.com/vtfr/bossanova/model"
	"github.com/vtfr/bossanova/service"
)

var _ = Describe("ExtractToken", func() {
	It("should extract a valid token from request and return it", func() {
		value, err := service.ExtractToken("Bearer token")
		Expect(value).To(Equal("token"))
		Expect(err).ToNot(HaveOccurred())
	})
	It("should return error if invalid token prefix", func() {
		value, err := service.ExtractToken("invalid token")
		Expect(value).To(Equal(""))
		Expect(err).To(HaveOccurred())
	})
})

var _ = Describe("Authentication", func() {
	var auth service.Authenticator
	var mockCtrl *gomock.Controller
	var mockStore *mocks.MockStore
	var token string

	sample := model.NewUser("sample", "password", "admin")

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockStore = mocks.NewMockStore(mockCtrl)
		auth = service.NewAuthenticator(mockStore, []byte("secret"))
	})
	AfterEach(func() {
		mockCtrl.Finish()
	})
	It("should generate valid tokens", func() {
		var err error
		token, err = auth.CreateToken(sample)
		Expect(err).ToNot(HaveOccurred())
	})
	It("should fail to parse invalid tokens", func() {
		_, err := auth.AuthenticateToken("invalid")
		Expect(err).To(HaveOccurred())
	})
	It("should be able to parse valid tokens if user exists", func() {
		mockStore.EXPECT().GetUser(sample.Username).
			Return(sample, nil).
			Times(1)

		user, err := auth.AuthenticateToken(token)
		Expect(err).ToNot(HaveOccurred())
		Expect(user).To(Equal(sample))
	})
	It("should return error if no users exists", func() {
		mockStore.EXPECT().GetUser(sample.Username).
			Return(nil, errors.New("no user")).
			Times(1)

		_, err := auth.AuthenticateToken(token)
		Expect(err).To(HaveOccurred())
	})
})
