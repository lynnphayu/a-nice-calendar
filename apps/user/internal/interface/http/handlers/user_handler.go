package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/subscription-tracker/user/internal/core/application"
)

type UserHandler struct {
	userService *application.UserService
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type PasswordlessLoginRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type VerifyTokenRequest struct {
	Token string `json:"token" binding:"required"`
}

type UpdateUserRequest struct {
	Name     string `json:"name"`
	Password string `json:"password,omitempty" binding:"omitempty,min=6"`
}

func NewUserHandler(userService *application.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	response, err := h.userService.Register(req.Email, req.Name, req.Password)
	if err == application.ErrUserAlreadyExists {
		c.JSON(http.StatusConflict, gin.H{"message": "user already exists"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to register user"})
		return
	}

	c.JSON(http.StatusCreated, response)
}

func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	response, err := h.userService.Login(req.Email, req.Password)
	if err == application.ErrInvalidCredentials {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid credentials"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to login"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) GetUser(c *gin.Context) {
	userID := c.Param("id")
	user, err := h.userService.GetUserByID(userID)
	if err == application.ErrUserNotFound {
		c.JSON(http.StatusNotFound, gin.H{"message": "user not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to get user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	userID := c.Param("id")
	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err := h.userService.UpdateUser(userID, req.Name, req.Password)
	if err == application.ErrUserNotFound {
		c.JSON(http.StatusNotFound, gin.H{"message": "user not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user updated successfully"})
}

func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	userID := c.GetString("user_id")
	fmt.Println(userID)
	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to get users"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	err := h.userService.DeleteUser(userID)
	if err == application.ErrUserNotFound {
		c.JSON(http.StatusNotFound, gin.H{"message": "user not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
}

func (h *UserHandler) InitiatePasswordlessLogin(c *gin.Context) {
	var req PasswordlessLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err := h.userService.InitiatePasswordlessLogin(req.Email)
	if err == application.ErrUserNotFound {
		c.JSON(http.StatusNotFound, gin.H{"message": "user not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to initiate passwordless login", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "login link sent to email"})
}

func (h *UserHandler) VerifyLoginToken(c *gin.Context) {
	var req VerifyTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	response, err := h.userService.VerifyLoginToken(req.Token)
	if err == application.ErrInvalidToken {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid or expired token"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to verify token"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) InitiateLogin(c *gin.Context) {
	state := h.userService.GenerateOAuthState()
	c.SetCookie("oauth_state", state, 3600, "/", "", false, true)

	spotifyAuthURL := h.userService.GetSpotifyAuthURL(state)
	fmt.Println(spotifyAuthURL)
	c.Redirect(http.StatusTemporaryRedirect, spotifyAuthURL)
}

func (h *UserHandler) HandleCallback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")
	savedState, _ := c.Cookie("oauth_state")

	if state == "" || state != savedState {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid state parameter"})
		return
	}

	response, err := h.userService.HandleSpotifyCallback(code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to handle Spotify callback", "error": err.Error()})
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, "http://localhost:3000/login/callback/spotify?token="+response.Token)
}
