package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	handler "stepik.leoscode.http/app/handler"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()

	router.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "hello web",
		})
	})

	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	router.GET("/internal/v1/health", handler.Health)
	router.GET("/internal/v1/Truncate", handler.Truncate)
	router.POST("/api/v1/auth/login", handler.LoginMiddleware(), handler.Login)
	router.POST("/api/v1/threads", handler.CreateThreads)
	router.GET("/api/v1/threads", handler.ShowThreads)
	router.GET("/api/v1/threads/:id", handler.GetThread)
	router.DELETE("/api/v1/threads/:id", handler.DeleteThread)
	router.POST("/api/v1/threads/:id/posts", handler.CreatePosts)
	router.GET("/api/v1/threads/:id/posts", handler.GetPosts)
	router.PATCH("/api/v1/threads/:id/lock", handler.LockThread)
	router.PATCH("/api/v1/threads/:id", handler.PatchThread)
	router.DELETE("/api/v1/posts/:id", handler.DeletePost)
	router.PATCH("/api/v1/posts/:id", handler.PatchPost)
	router.PUT("/api/v1/threads/:id", handler.PutThread)
	router.POST("/api/v1/threads/:id/attachments", handler.UploadAttachment)
	router.GET("/api/v1/threads/:id/attachments", handler.ShowAttachments)
	router.GET("/api/v1/attachments/:id/file", handler.DownloadAttachment)
	router.DELETE("/api/v1/attachments/:id", handler.DeleteAttachment)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGABRT)

	go func() {
		fmt.Println("http://localhost:8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("server error", err)
			return
		}
	}()

	<-done

	ctx, close := context.WithTimeout(context.Background(), 5*time.Second)
	defer close()

	if err := server.Shutdown(ctx); err != nil {
		fmt.Println("Shutdown error", err)
	} else {
		fmt.Println("Server is stopped")
	}
}
