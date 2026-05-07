package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"stepik.leoscode.http/app/models"
	"stepik.leoscode.http/app/service"

	"github.com/gin-gonic/gin"
)

func CreateThreads(c *gin.Context) {
	userID := c.GetHeader("X-User-Id")
	idempotencyKey := c.GetHeader("X-Idempotency-Key")
	if userID == "" {
		c.JSON(401, gin.H{
			"code":    "unauthorized",
			"message": "missing user id",
		})
		return
	}
	if userID == "" {
		c.JSON(401, gin.H{
			"code":    "unauthorized",
			"message": "missing idempotency key",
		})
		return
	}

	var thread models.Thread

	if err := c.BindJSON(&thread); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Wrong JSON",
		})
		return
	}

	thread, status, err := service.CreateThread(userID, idempotencyKey, thread)
	if err != nil {
		if status == 409 {
			c.JSON(409, gin.H{
				"code": "conflict",
			})
			return
		}

		c.JSON(500, gin.H{
			"code": fmt.Sprint("internal_error", err),
		})
		return
	}

	c.JSON(status, thread)
}

func ShowThreads(c *gin.Context) {
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
	param.Tag = c.DefaultQuery("tag", "")
	param.Author = c.DefaultQuery("author_id", "")

	threads, total, err := service.GetThreads(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Items": threads,
		"meta": gin.H{
			"limit":  param.Limit,
			"offset": param.Offset,
			"total":  total,
		},
	})
}
