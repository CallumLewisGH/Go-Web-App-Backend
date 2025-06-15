package api

import (
	"fmt"
	"os"

	_ "github.com/CallumLewisGH/Generic-Service-Base/docs"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
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
	router := gin.Default() //Change to prod

	// Session configuration
	store := cookie.NewStore(
		[]byte(os.Getenv("GIN_SESSION_SECRET")),
	)

	// Session middleware
	router.Use(sessions.Sessions("mysession", store))

	// Swagger setup
	router.GET("/swagger/*any", swagger.WrapHandler(swaggerFiles.Handler))

	s := &Server{
		Engine: router,
	}

	fmt.Println("Server Creation Succeeded")
	return s
}
