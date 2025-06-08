package routes

import (
	"net/http"

	_ "github.com/CallumLewisGH/Generic-Service-Base/docs"
	"github.com/CallumLewisGH/Generic-Service-Base/internal/api"
	command "github.com/CallumLewisGH/Generic-Service-Base/internal/api/handlers/commands"
	query "github.com/CallumLewisGH/Generic-Service-Base/internal/api/handlers/queries"
	userModel "github.com/CallumLewisGH/Generic-Service-Base/internal/domain/user"
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
		databaseGroup.POST("", createUser)
	}
}

// Get users godoc
// @Summary Returns all the users
// @Description Gets all the users
// @Tags users
// @Produce json
// @Success 200 {object} userModel.UserDTO "Returns a pagenated list of users"
// @Router /users [get]
func getUsers(c *gin.Context) {
	user, err := query.GetAllUsersQuery()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}

	c.JSON(http.StatusOK, user)
}

// Create user godoc
// @Summary Creates a new user
// @Description Creates a new user with the provided details
// @Tags users
// @Accept json
// @Produce json
// @Param user body userModel.UserRequest true "User details"
// @Success 201 {object} userModel.UserDTO "Returns the created user"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users [post]
func createUser(c *gin.Context) {
	var req userModel.UserRequest

	// Bind JSON request to UserRequest struct
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Validate required fields
	if req.Username == "" || req.Email == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username, email, and password are required"})
		return
	}

	// Call command to create user
	user, err := command.CreateUserCommand(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, user)
}
