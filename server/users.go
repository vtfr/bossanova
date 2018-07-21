package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vtfr/bossanova/model"
)

type createUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required"`
}

type updateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func getUsers(c *gin.Context) {
	users, err := GetStore(c).AllUsers()
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}

func getUser(c *gin.Context) {
	user, err := GetStore(c).GetUser(c.Param("username"))
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

func createUser(c *gin.Context) {
	var data createUserRequest

	if err := c.BindJSON(&data); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	user := &model.User{
		Username:       data.Username,
		HashedPassword: model.HashPassword(data.Password),
		Role:           data.Role,
	}

	if err := model.Validate(user); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if err := GetStore(c).CreateUser(user); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusCreated)
}

func updateUser(c *gin.Context) {
	var data updateUserRequest

	// Bind request data to object
	if err := c.BindJSON(&data); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	user := &model.User{
		Username: data.Username,
		Role:     data.Role,
	}

	if data.Password != "" {
		user.SetPassword(data.Password)
	}

	if err := model.Validate(user); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if err := GetStore(c).UpdateUser(user); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func deleteUser(c *gin.Context) {
	if err := GetStore(c).DeleteUser(c.Param("username")); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusNoContent)
}
