package handler

import (
	"net/http"
	"path/filepath"
	"strconv"

	"stepik.leoscode.http/app/models"
	"stepik.leoscode.http/app/service"

	"github.com/gin-gonic/gin"
)

func UploadAttachment(c *gin.Context) {
	var attc models.Attachments

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Uncorect data",
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

	userID := c.GetHeader("X-User-Id")
	attc.Filename = file.Filename
	attc.MimeType = file.Header.Get("Content-Type")
	attc.ThreadId = id
	attc.Size = int(file.Size)

	dst := filepath.Join("./test_Files/", filepath.Base(file.Filename))
	attc.Path = dst
	c.SaveUploadedFile(file, dst)

	attc, err = service.UploudFile(attc, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, attc)
}

func ShowAttachments(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Wrong id",
		})
		return
	}

	attcs, err := service.ShowFiles(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": attcs,
	})
}

func DownloadAttachment(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "No id",
		})
		return
	}

	filePath, fileName, err := service.GetFile(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	c.FileAttachment(filePath, fileName)
}

func DeleteAttachment(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "No id",
		})
		return
	}

	if err := service.DeleteFile(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusNoContent, "No Content")
}
