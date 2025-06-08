package routes

import (
	"net/http"

	_ "github.com/CallumLewisGH/Generic-Service-Base/docs"
	"github.com/CallumLewisGH/Generic-Service-Base/internal/api"
	query "github.com/CallumLewisGH/Generic-Service-Base/internal/api/handlers/queries"
	_ "github.com/CallumLewisGH/Generic-Service-Base/internal/domain/models"
	"github.com/gin-gonic/gin"
)

// Register User Routes godoc
// @Summary User-related endpoints
// @Description User-related endpoints
// @Tags User
func RegisterUserRoutes(s *api.Server) {
	databaseGroup := s.Group("/users")
	{
		databaseGroup.GET("", getUsers)
	}
}

// Get users godoc
// @Summary Returns all the users
// @Description Gets all the users
// @Tags users
// @Produce json
// @Success 200 {object} models.User "Returns a pagenated list of users"
// @Router /users [get]
func getUsers(c *gin.Context) {
	result, err := query.GetAllUsersQuery()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}

	c.JSON(http.StatusOK, result)
}

// Param id path string true "User ID" Format(uuid) Example(550e8400-e29b-41d4-a716-446655440000)
