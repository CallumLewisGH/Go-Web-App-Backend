package routes

import (
	"net/http"

	_ "github.com/CallumLewisGH/Generic-Service-Base/docs"
	"github.com/CallumLewisGH/Generic-Service-Base/internal/api"
	command "github.com/CallumLewisGH/Generic-Service-Base/internal/api/handlers/commands"
	query "github.com/CallumLewisGH/Generic-Service-Base/internal/api/handlers/queries"
	"github.com/gin-gonic/gin"
)

// RegisterDatabaseRoutes godoc
// @Summary Database-related endpoints
// @Description Operations for database management
// @Tags database
func RegisterDatabaseRoutes(s *api.Server) {
	databaseGroup := s.Group("/database")
	{
		databaseGroup.GET("/health", getHealth)
		databaseGroup.POST("/migrations", runMigrations)
	}
}

// Get database health godoc
// @Summary Returns the state of the connected database
// @Description Gets the state of the connected database
// @Tags database
// @Produce json
// @Success 200 {string} string "Returns database health status"
// @Router /database/health [get]
func getHealth(c *gin.Context) {
	health := query.CheckDatabaseHealth()
	c.JSON(http.StatusOK, health)
}

// Post database migrations godoc
// @Summary Runs the migrations from GORM
// @Description Runs the database migrations from GORM
// @Tags database
// @Produce json
// @Success 200 {string} string "Retuns success message"
// @Failure 500 {object} error "Retuns any errors"
// @Router /database/migrations [post]
func runMigrations(c *gin.Context) {
	migrations := command.RunDatabaseMigrations()
	if migrations != nil {
		c.JSON(http.StatusInternalServerError, migrations)
	}
	c.JSON(http.StatusOK, "Migrations Successfull")
}
