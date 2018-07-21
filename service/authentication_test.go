package service_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/franela/goblin"
	"github.com/golang/mock/gomock"
	"github.com/vtfr/bossanova/mocks"
	"github.com/vtfr/bossanova/model"
	"github.com/vtfr/bossanova/service"
)

func TestExtractToken(t *testing.T) {
	g := goblin.Goblin(t)
	g.Describe("ExtractToken", func() {
		g.It("Should extract a valid token from request and return it", func() {
			req, err := http.NewRequest(http.MethodGet, "/", nil)
			if err != nil {
				t.Fail()
			}

			// adds authorization header
			req.Header.Set("Authorization", "Bearer token")

			value, err := service.ExtractToken(req)
			g.Assert(err).Equal(nil)
			g.Assert(value).Equal("token")
		})
		g.It("Should return error if no token", func() {
			req, err := http.NewRequest(http.MethodGet, "/", nil)
			if err != nil {
				t.Fail()
			}

			_, err = service.ExtractToken(req)
			g.Assert(err != nil).IsTrue()
		})
		g.It("Should return error if invalid token prefix", func() {
			req, err := http.NewRequest(http.MethodGet, "/", nil)
			if err != nil {
				t.Fail()
			}

			// adds authorization header
			req.Header.Set("Authorization", "token")

			_, err = service.ExtractToken(req)
			g.Assert(err != nil).IsTrue()
		})
	})

}

func TestAuthenticator(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	userStore := mocks.NewMockStore(mockCtrl)
	auth := service.NewAuthenticator(userStore, []byte("secret"))

	sample := &model.User{
		Username: "sample",
		Role:     "admin",
	}

	var token string

	g := goblin.Goblin(t)
	g.Describe("Authenticator", func() {
		g.It("Should generate valid tokens", func() {
			// generate token
			var err error
			token, err = auth.CreateToken(sample)
			g.Assert(err == nil).IsTrue("failed generating token")
		})
		g.It("Should fail to parse invalid tokens", func() {
			user, err := auth.AuthenticateToken("invalid")
			g.Assert(user == nil).IsTrue()
			g.Assert(err != nil).IsTrue()
		})
		g.It("Should be able to parse valid tokens if user exists", func() {
			userStore.EXPECT().GetUser(sample.Username).
				Return(sample, nil).
				Times(1)

			user, err := auth.AuthenticateToken(token)
			g.Assert(err == nil).IsTrue("failed getting/parsing user")

			// verify same users
			g.Assert(user.Username == sample.Username)
		})
		g.It("Should return error if no users exists", func() {
			userStore.EXPECT().GetUser(sample.Username).
				Return(nil, errors.New("no user")).
				Times(1)

			// parse request
			user, err := auth.AuthenticateToken(token)
			g.Assert(err != nil).IsTrue()
			g.Assert(user == nil).IsTrue()
		})
	})
}
