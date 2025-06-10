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

func googleLogin(c *gin.Context) {
	session := sessions.Default(c)
	state, _ := authentication.GenerateState()
	session.Set("oauth_state", state)
	session.Save()

	url := authentication.Config.AuthCodeURL(state)
	c.Redirect(http.StatusFound, url)
}

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

func getAuthenticatedUser(c *gin.Context) {
	userAuthId := c.MustGet("userAuthId").(string)

	user, err := query.GetUserByAuthIdQuery(userAuthId)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func logout(c *gin.Context) {
	c.SetCookie("jwt", "", -1, "/", os.Getenv("DOMAIN"), true, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out"})
}
