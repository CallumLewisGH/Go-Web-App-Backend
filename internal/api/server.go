package api

import (
	"fmt"
	"os"
	"time"

	_ "github.com/CallumLewisGH/Generic-Service-Base/docs"
	"github.com/CallumLewisGH/Generic-Service-Base/internal/api/middleware"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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

	_ = godotenv.Load()

	// Global middleware middleware application
	router.Use(sessions.Sessions("mysession", store))
	router.Use(middleware.NewRateLimiter(5, time.Second, os.Getenv("REDIS_URL")))

	// Swagger setup
	router.GET("/swagger/*any", swagger.WrapHandler(swaggerFiles.Handler))

	s := &Server{
		Engine: router,
	}

	fmt.Println("Server Creation Succeeded")
	return s
}
