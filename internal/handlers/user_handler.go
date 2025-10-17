package handlers

import (
	"net/http"
	"shop/internal/database"
	"shop/internal/services"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// @Summary Register a new user
// @Description Create a new user account with email, username, and password.
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body services.RegisterRequest true "User registration data"
// @Success 201 {object} services.RegisterResponse "Created user"
// @Failure 400 {object} map[string]string "Invalid input (validation error)"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/register [post]
func (r *Router) Register(c *gin.Context) {
	var register services.RegisterRequest

	if err := c.ShouldBindJSON(&register); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(register.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}
	register.Password = string(hashedPassword)

	user := database.User{
		Username: register.Username,
		Password: string(hashedPassword),
		Email:    register.Email,
	}
	err = r.models.Users.Insert(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't create user"})
		return
	}

	token, err := r.getJwtToken(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	c.JSON(http.StatusCreated, services.RegisterResponse{Email: register.Email, Username: register.Username, Token: token})
}

// login authenticates a user and returns a JWT token
// @Summary Login user
// @Description Authenticate user by username and password and return a JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param credentials body services.LoginRequest true "User login credentials"
// @Success 200 {object} services.LoginResponse "JWT token"
// @Failure 400 {object} map[string]string "Invalid input (validation error)"
// @Failure 401 {object} map[string]string "Invalid username or password"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/login [post]
func (r *Router) Login(c *gin.Context) {
	var login services.LoginRequest

	if err := c.ShouldBindJSON(&login); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existingUser, err := r.models.Users.FindByUsername(login.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(login.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	token, err := r.getJwtToken(existingUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	c.JSON(http.StatusOK, services.LoginResponse{Token: token})
}

func (r *Router) getJwtToken(user *database.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": user.ID,
		"exp":    time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString([]byte(r.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
