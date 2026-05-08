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
	if idempotencyKey == "" {
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
	param.Sort = c.DefaultQuery("sort", "old")
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

func GetThread(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Wrong id",
		})
		return
	}

	thread, err := service.GetThreadById(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, thread)
}

func DeleteThread(c *gin.Context) {
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

	if err := service.DeleteThread(userID, id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusNoContent, "No Content")
}

func PatchThread(c *gin.Context) {
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

	var newThread models.Thread
	if err := c.BindJSON(&newThread); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Wrong JSON",
		})
		return
	}
	thread, err := service.UpdateThread(newThread, userID, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id":         thread.Id,
		"title":      thread.Title,
		"content":    thread.Content,
		"updated_at": thread.UpdatedAt,
	})
}

func LockThread(c *gin.Context) {
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

	var LockThread models.Thread
	if err := c.BindJSON(&LockThread); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Wrong JSON",
		})
		return
	}

	if err := service.LockThread(LockThread.IsLocked, userID, id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":        id,
		"is_locked": LockThread.IsLocked,
	})
}

func SortThreads(c *gin.Context) {
	sort := c.DefaultQuery("sort", "old")

	threds, err := service.GetSortThreads(sort)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"Threads": threds,
	})
}
