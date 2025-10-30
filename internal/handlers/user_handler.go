package handlers

import (
	"net/http"
	"shop/internal/dto"
	"shop/internal/services"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *services.UserService
}

// Register user
// @Summary Registers a new user
// @Description Create a new user account with email, username, and password.
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body dto.RegisterRequest true "User registration data"
// @Success 201 {object} dto.RegisterResponse "Created user"
// @Failure 400 {object} map[string]string "Invalid input (validation error)"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/auth/register [post]
func (uh *UserHandler) Register(c *gin.Context) {
	var register dto.RegisterRequest

	if err := c.ShouldBindJSON(&register); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	regRes, errMsg, status := uh.userService.Register(register)
	if errMsg != "" {
		c.AbortWithStatusJSON(status, gin.H{"error": errMsg})
		return
	}

	c.JSON(status, regRes)
}

// Login authenticates a user and returns a JWT token
// @Summary Login user
// @Description Authenticate user by username and password and return a JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param credentials body dto.LoginRequest true "User login credentials"
// @Success 200 {object} dto.LoginResponse "JWT token"
// @Failure 400 {object} map[string]string "Invalid input (validation error)"
// @Failure 401 {object} map[string]string "Invalid username or password"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/auth/login [post]
func (uh *UserHandler) Login(c *gin.Context) {
	var login dto.LoginRequest

	if err := c.ShouldBindJSON(&login); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logRes, errMsg, status := uh.userService.Login(login)
	if errMsg != "" {
		c.AbortWithStatusJSON(status, gin.H{"error": errMsg})
		return
	}

	c.JSON(status, logRes)
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}
