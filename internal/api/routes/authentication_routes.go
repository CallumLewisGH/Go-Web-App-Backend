package routes

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/CallumLewisGH/Generic-Service-Base/internal/api"
	"github.com/CallumLewisGH/Generic-Service-Base/internal/api/authentication"
	command "github.com/CallumLewisGH/Generic-Service-Base/internal/api/handlers/commands"
	query "github.com/CallumLewisGH/Generic-Service-Base/internal/api/handlers/queries"
	"github.com/CallumLewisGH/Generic-Service-Base/internal/api/middleware"
	userModel "github.com/CallumLewisGH/Generic-Service-Base/internal/domain/user"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// RegisterAuthenticationRoutes godoc
// @Summary Authentication-related endpoints
// @Description Operations for authentication management
// @Tags authentication
func RegisterAuthenticationRoutes(s *api.Server) {
	authenticationGroup := s.Group("/authentication")
	{
		authenticationGroup.GET("/login", googleLogin)
		authenticationGroup.POST("/callback", googleCallback)
		authenticationGroup.POST("/logout", logout)
		authenticationGroup.GET("/user", getAuthenticatedUser)
	}
}

// googleLogin godoc
// @Summary Initiate Google OAuth2 login
// @Description Redirects to Google's OAuth2 consent page
// @Tags authentication
// @Produce json
// @Success 302 {string} string "Redirects to Google OAuth"
// @Failure 500 {object} error "Returns error if session state cannot be generated"
// @Router /authentication/login [get]
func googleLogin(c *gin.Context) {
	session := sessions.Default(c)
	state, _ := authentication.GenerateState()
	session.Set("oauth_state", state)
	session.Save()

	url := authentication.Config.AuthCodeURL(state)
	c.Redirect(http.StatusFound, url)
}

// googleCallback godoc
// @Summary Google OAuth2 callback
// @Description Handles the callback from Google OAuth2, creates user if new, and sets JWT cookie
// @Tags authentication
// @Accept json
// @Produce json
// @Param state query string true "OAuth state token"
// @Param code query string true "OAuth authorization code"
// @Success 302 {string} string "Redirects to /home on success"
// @Failure 400 {object} error "Invalid state token"
// @Failure 500 {object} error "Returns error if token exchange, user info fetch, or JWT generation fails"
// @Failure 409 {object} error "Returns error if username derived from email is already in use"
// @Router /authentication/callback [post]

func googleCallback(c *gin.Context) {
	session := sessions.Default(c)

	if c.Query("state") != session.Get("oauth_state") {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid state token"})
		return
	}

	token, err := authentication.Config.Exchange(c, c.Query("code"))

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
		return
	}

	client := authentication.Config.Client(c, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}
	defer resp.Body.Close()

	var googleUser struct {
		Email  string `json:"email"`
		AuthId string `json:"id"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode google's user info"})
		return
	}

	_, err = query.GetUserByAuthIdQuery(googleUser.AuthId)

	//If existing user does not exist create one
	if err != nil {
		userReq := userModel.CreateUserRequest{
			Email:    googleUser.Email,
			AuthId:   googleUser.AuthId,
			Username: strings.Split(googleUser.Email, "@")[0],
		}

		_, err := command.CreateUserCommand(userReq)

		if err != nil {
			c.JSON(http.StatusConflict, gin.H{"error": "The username based off of the email: " + userReq.Email + " is already in use." + err.Error()})
			return
		}
	}

	jwtToken, err := middleware.GenerateJWT(googleUser.AuthId, googleUser.Email)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWT"})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("jwt", jwtToken, int(time.Hour*24*7), "/", os.Getenv("DOMAIN"), true, true)
	c.Redirect(http.StatusFound, "/home")
}

// getAuthenticatedUser godoc
// @Summary Get authenticated user info
// @Description Returns information about the currently authenticated user
// @Tags authentication
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} userModel.UserDTO "Returns user object"
// @Failure 404 {object} error "User not found"
// @Router /authentication/user [get]
func getAuthenticatedUser(c *gin.Context) {
	userAuthId := c.MustGet("userAuthId").(string)

	user, err := query.GetUserByAuthIdQuery(userAuthId)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// logout godoc
// @Summary Logout current user
// @Description Clears the JWT cookie, effectively logging out the user
// @Tags authentication
// @Produce json
// @Success 200 {object} object "Returns success message"
// @Router /authentication/logout [post]

func logout(c *gin.Context) {
	c.SetCookie("jwt", "", -1, "/", os.Getenv("DOMAIN"), true, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out"})
}
