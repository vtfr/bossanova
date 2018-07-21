package server_test

import (
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/franela/goblin"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/vtfr/bossanova/mocks"
	"github.com/vtfr/bossanova/server"
)

func TestMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	g := goblin.Goblin(t)
	g.Describe("StoreMiddleware", func() {
		g.It("Should store a cloned store in a context", func() {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			storeMock := mocks.NewMockStore(mockCtrl)
			storeMock.EXPECT().Clone().Return(storeMock).Times(1)
			storeMock.EXPECT().Close().Times(1)

			r := gin.New()
			r.Use(server.StoreMiddleware(storeMock))
			r.GET("/", func(c *gin.Context) {
				defer func() {
					if err := recover(); err != nil {
						g.Fail(err)
					}
				}()

				// fails if no store is found
				server.GetStore(c)
			})

			req, _ := http.NewRequest(http.MethodGet, "/", nil)
			r.ServeHTTP(httptest.NewRecorder(), req)
		})
	})
}
