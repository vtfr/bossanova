package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vtfr/bossanova/model"
)

func getPost(c *gin.Context) {
	post, err := GetStore(c).GetPost(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"post": post,
	})
}

func updatePost(c *gin.Context) {
	var post model.Post
	if err := c.BindJSON(&post); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	post.ID = c.Param("id")
	if err := GetStore(c).UpdatePost(&post); err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func deletePost(c *gin.Context) {
	if err := GetStore(c).DeletePost(c.Param("id")); err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func getThread(c *gin.Context) {
	thread, err := GetStore(c).GetThread(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"thread": thread,
	})
}
