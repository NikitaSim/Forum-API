package handler

import (
	"net/http"
	"regexp"

	"stepik.leoscode.http/app/models"
	"stepik.leoscode.http/app/service"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	req := c.MustGet("user").(models.Users)
	id, err := service.LoginUser(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Error adding user",
		})
	}

	c.JSON(http.StatusCreated, gin.H{
		"user_id": id,
	})
}

func LoginMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.Users
		tmpl := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
		if err := c.ShouldBind(&user); err != nil || !tmpl.MatchString(user.Name) { // добавить regexp
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": "Incorrect data",
			})
		}

		c.Set("user", user)
		c.Next()
	}
}
