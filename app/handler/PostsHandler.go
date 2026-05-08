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

	c.JSON(http.StatusCreated, post)
}

func GetPosts(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Wrong id",
		})
		return
	}

	var param models.Parameters
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Wrong limit",
		})
		return
	}
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Wrong offset",
		})
		return
	}
	param.Limit = limit
	param.Offset = offset

	posts, total, err := service.GetPosts(id, param)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Items": posts,
		"meta": gin.H{
			"limit":  param.Limit,
			"offset": param.Offset,
			"total":  total,
		},
	})
}

func DeletePost(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Wrong id",
		})
		return
	}
	userID := c.GetHeader("X-User-Id")
	if userID == "" {
		c.JSON(401, gin.H{
			"code":    "unauthorized",
			"message": "missing user id",
		})
		return
	}

	if err := service.DeletePost(userID, id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusNoContent, "No Content")
}

func PatchPost(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Wrong id",
		})
		return
	}
	userID := c.GetHeader("X-User-Id")
	if userID == "" {
		c.JSON(401, gin.H{
			"code":    "unauthorized",
			"message": "missing user id",
		})
		return
	}

	var cont models.Contents
	if err := c.BindJSON(&cont); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Wrong JSON",
		})
	}

	update, err := service.UpdatePost(userID, id, cont.Content)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         id,
		"content":    cont.Content,
		"updated_at": *update,
	})
}
