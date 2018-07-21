package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vtfr/bossanova/service"
)

// post create a new post
func post(auth service.Authorizator) gin.HandlerFunc {
	return func(c *gin.Context) {
		var data service.CreatePostData
		if err := c.Bind(&data); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		data.IP = c.ClientIP()
		post, err := service.CreatePost(auth, GetStore(c), GetUser(c), &data)
		if err != nil {
			c.Error(err)
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"post": post,
		})
	}
}
