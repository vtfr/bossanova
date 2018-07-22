package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vtfr/bossanova/common"
	"github.com/vtfr/bossanova/model"
)

func getBoards(c *gin.Context) {
	boards, err := GetStore(c).AllBoards()
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"boards": boards,
	})
}

func getBoard(c *gin.Context) {
	board, err := GetStore(c).GetBoard(c.Param("uri"))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"board": board,
	})
}

func createBoard(c *gin.Context) {
	var board *model.Board
	if err := c.BindJSON(&board); err != nil {
		c.Error(common.ErrBadRequest)
		return
	}

	if err := model.Validate(board); err != nil {
		c.Error(err)
		return
	}

	if err := GetStore(c).CreateBoard(board); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusCreated)
}

func updateBoard(c *gin.Context) {
	var board *model.Board
	if err := c.BindJSON(&board); err != nil {
		c.Error(common.ErrBadRequest)
		return
	}

	board.URI = c.Param("uri")

	if err := model.Validate(board); err != nil {
		c.Error(err)
		return
	}

	if err := GetStore(c).UpdateBoard(board); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}

func deleteBoard(c *gin.Context) {
	if err := GetStore(c).DeleteBoard(c.Param("uri")); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
