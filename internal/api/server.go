package api

import (
	"fmt"

	_ "github.com/CallumLewisGH/Generic-Service-Base/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	swagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	*gin.Engine
}

// @title Generic-Service-Base API
// @version 1.0
func NewServer() *Server {
	fmt.Println("Creating Server...")

	// Create Gin instance (with logger and recovery middleware)
	router := gin.Default()

	// Swagger setup
	router.GET("/swagger/*any", swagger.WrapHandler(swaggerFiles.Handler))

	s := &Server{
		Engine: router,
	}

	fmt.Println("Server Creation Succeeded")
	return s
}
