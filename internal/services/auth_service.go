package services

import (
	"context"
	"net/http"
	"shop/internal/database"
	"shop/internal/dto"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	models    *database.Models
	jwtSecret string
}

func (as *AuthService) Register(req dto.RegisterRequest) (resp *dto.RegisterResponse, errMsg string, status int) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "Something went wrong", http.StatusInternalServerError
	}
	req.Password = string(hashedPassword)

	user := database.User{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
		Type:     req.Type,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err = as.models.Users.Insert(&user, ctx)
	if err != nil {
		return nil, "Something went wrong", http.StatusInternalServerError
	}
	token, err := as.getJwtToken(&user)
	if err != nil {
		return nil, "Something went wrong", http.StatusInternalServerError
	}
	return &dto.RegisterResponse{Username: user.Username, Email: user.Email, Token: token}, "", http.StatusCreated
}

func (as *AuthService) Login(req dto.LoginRequest) (resp *dto.LoginResponse, errMsg string, status int) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	user, err := as.models.Users.FindByUsername(req.Username, ctx)
	if err != nil {
		return nil, "Invalid username or password", http.StatusUnauthorized
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, "Invalid username or password", http.StatusUnauthorized
	}

	token, err := as.getJwtToken(user)
	if err != nil {
		return nil, "Something went wrong", http.StatusInternalServerError
	}

	return &dto.LoginResponse{Token: token}, "", http.StatusOK
}

func (as *AuthService) GetUserFromContext(c *gin.Context) *database.User {
	cUser, exist := c.Get("user")
	if !exist {
		return nil
	}
	user, ok := cUser.(*database.User)
	if !ok {
		return nil
	}

	return user
}

func (as *AuthService) getJwtToken(user *database.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": user.ID,
		"exp":    time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString([]byte(as.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
