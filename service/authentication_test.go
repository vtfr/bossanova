package service_test

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/franela/goblin"
	"github.com/golang/mock/gomock"
	"github.com/vtfr/bossanova/mocks"
	"github.com/vtfr/bossanova/model"
	"github.com/vtfr/bossanova/service"
)

func createRequestWithToken(token string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	return req
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
			user, err := auth.GetUserFromRequest(createRequestWithToken("invalid"))
			g.Assert(user == nil).IsTrue()
			g.Assert(err != nil).IsTrue()
		})
		g.It("Should be able to parse valid tokens if user exists", func() {
			userStore.EXPECT().GetUser(sample.Username).
				Return(sample, nil).
				Times(1)

			user, err := auth.GetUserFromRequest(createRequestWithToken(token))
			g.Assert(err == nil).IsTrue("failed getting/parsing user")

			// verify same users
			g.Assert(user.Username == sample.Username)
		})
		g.It("Should return error if no users exists", func() {
			userStore.EXPECT().GetUser(sample.Username).
				Return(nil, errors.New("no user")).
				Times(1)

			// parse request
			user, err := auth.GetUserFromRequest(createRequestWithToken(token))
			g.Assert(err != nil).IsTrue()
			g.Assert(user == nil).IsTrue()
		})
	})
}
