package handler

import (
	"net/http"
	"strconv"

	"stepik.leoscode.http/app/models"
	"stepik.leoscode.http/app/service"

	"github.com/gin-gonic/gin"
)

func CreatePosts(c *gin.Context) {
	userID := c.GetHeader("X-User-Id")
	if userID == "" {
		c.JSON(401, gin.H{
			"code":    http.StatusUnauthorized,
			"message": "missing user id",
		})
		return
	}
	var content models.Contents
	if err := c.BindJSON(&content); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Wrong JSON",
		})
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Wrong id",
		})
		return
	}

	post, err := service.CreatePost(userID, id, content.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
		})
	}

	c.JSON(http.StatusOK, post)

}
