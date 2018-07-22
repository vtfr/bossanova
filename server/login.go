package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vtfr/bossanova/service"
)

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func login(auth service.Authenticator) gin.HandlerFunc {
	return func(c *gin.Context) {
		var login loginRequest
		if err := c.BindJSON(&login); err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		user, err := GetStore(c).GetUser(login.Username)
		if err != nil {
			c.Status(http.StatusForbidden)
			return
		}

		// verifies if right password
		if user.VerifyPassword(login.Password) {
			token, _ := auth.CreateToken(user)
			c.JSON(http.StatusOK, gin.H{
				"token": token,
			})
		} else {
			c.Status(http.StatusForbidden)
		}
	}
}
