package handler

import (
	"net/http"

	"stepik.leoscode.http/app/service"

	"github.com/gin-gonic/gin"
)

func Health(c *gin.Context) {
	header := c.GetHeader("X-Request-Id")
	if header == "" {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":       http.StatusOK,
			"X-Request-Id": header,
		})
	}
}

func Truncate(c *gin.Context) {
	err := service.ClearTables()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Еable deletion error",
		})
	}

	c.Status(http.StatusNoContent)
}
