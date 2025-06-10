package routes

import (
	"net/http"

	_ "github.com/CallumLewisGH/Generic-Service-Base/docs"
	"github.com/CallumLewisGH/Generic-Service-Base/internal/api"
	command "github.com/CallumLewisGH/Generic-Service-Base/internal/api/handlers/commands"
	query "github.com/CallumLewisGH/Generic-Service-Base/internal/api/handlers/queries"
	userModel "github.com/CallumLewisGH/Generic-Service-Base/internal/domain/user"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Register User Routes godoc
// @Summary User-related endpoints
// @Description User-related endpoints
// @Tags User
func RegisterUserRoutes(s *api.Server) {
	routeGroup := s.Group("/users")
	{
		routeGroup.GET("", getUsers)
		routeGroup.GET("/id", getUserById)
		routeGroup.POST("", createUser)
		routeGroup.PUT("/id", updateUserById)
		routeGroup.DELETE("/id", deleteUserById)
	}
}

// Get users godoc
// @Summary Returns all the users
// @Description Gets all the users
// @Tags users
// @Produce json
// @Success 200 {object} userModel.UserDTO "Returns a pagenated list of users"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users [get]
func getUsers(c *gin.Context) {
	user, err := query.GetAllUsersQuery()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// Get user by ID godoc
// @Summary Returns the user with the specified ID
// @Description Gets the user where ID is passed in the user_id header
// @Tags users
// @Produce json
// @Param user_id header string true "User ID to retrieve"
// @Success 200 {object} userModel.UserDTO "Returns the requested user"
// @Failure 400 {object} map[string]string "Missing or invalid user ID header"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users/id [get]
func getUserById(c *gin.Context) {
	// Get user ID from custom header
	userIDStr := c.GetHeader("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id header is required"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID must be a valid GUID"})
		return
	}

	user, err := query.GetUserByIdQuery(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// Create user godoc
// @Summary Creates a new user
// @Description Creates a new user with the provided details
// @Tags users
// @Accept json
// @Produce json
// @Param user body userModel.CreateUserRequest true "User details"
// @Success 201 {object} userModel.UserDTO "Returns the created user"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users [post]
func createUser(c *gin.Context) {
	var req userModel.CreateUserRequest

	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if req.Username == "" || req.Email == "" || req.AuthId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username, email, and authId are required"})
		return
	}

	user, err := command.CreateUserCommand(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// UpdateUser godoc
// @Summary Updates the user with the specified ID
// @Description Updates the user where ID is passed in the user_id header with the provided details
// @Tags users
// @Accept json
// @Produce json
// @Param user_id header string true "User ID to update"
// @Param user body userModel.UpdateUserRequest true "User update details"
// @Success 200 {object} userModel.UserDTO "Successfully updated user"
// @Failure 400 {object} map[string]string "Missing or invalid user ID header or invalid request body"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users/id [put]
func updateUserById(c *gin.Context) {
	userIDStr := c.GetHeader("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id header is required"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID must be a valid GUID"})
		return
	}

	var user userModel.UpdateUserRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	updatedUser, err := command.UpdateUserByIdCommand(userID, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}

// DeleteUser godoc
// @Summary Deletes the user with the specified ID
// @Description Deletes the user where ID is passed in the user_id header
// @Tags users
// @Produce json
// @Param user_id header string true "User ID to delete"
// @Success 200 {object} userModel.UserDTO "Successfull deletion returns the model of the deleted user"
// @Failure 400 {object} map[string]string "Missing or invalid user ID header"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users/id [delete]
func deleteUserById(c *gin.Context) {
	// Get user ID from custom header
	userIDStr := c.GetHeader("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id header is required"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID must be a valid GUID"})
		return
	}

	user, err := command.DeleteUserByIdCommand(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}
