package services

import (
	"context"
	"net/http"
	"shop/internal/dto"
	"shop/internal/models"
	"shop/internal/repositories"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepository repositories.UserRepository
	cartRepository repositories.CartRepository
	jwtSecret      string
}

func (us *UserService) Register(req dto.RegisterRequest) (resp *dto.RegisterResponse, errMsg string, status int) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "Something went wrong", http.StatusInternalServerError
	}
	req.Password = string(hashedPassword)

	user := models.User{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
		Type:     req.Type,
	}
	userCtx, userCancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer userCancel()
	err = us.userRepository.Insert(&user, userCtx)
	if err != nil {
		return nil, "Something went wrong", http.StatusInternalServerError
	}

	if user.Type == string(dto.TypeCustomer) {
		cart := models.Cart{
			UserID: user.ID,
			User:   &user,
		}
		cartCtx, cartCancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cartCancel()
		if err := us.cartRepository.Create(&cart, cartCtx); err != nil {
			return nil, "Something went wrong", http.StatusInternalServerError
		}
	}

	token, err := us.getJwtToken(&user)
	if err != nil {
		return nil, "Something went wrong", http.StatusInternalServerError
	}
	return &dto.RegisterResponse{Username: user.Username, Email: user.Email, Token: token}, "", http.StatusCreated
}

func (us *UserService) Login(req dto.LoginRequest) (resp *dto.LoginResponse, errMsg string, status int) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	user, err := us.userRepository.FindByEmail(req.Email, ctx)
	if err != nil {
		return nil, "Invalid username or password", http.StatusUnauthorized
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, "Invalid username or password", http.StatusUnauthorized
	}

	token, err := us.getJwtToken(user)
	if err != nil {
		return nil, "Something went wrong", http.StatusInternalServerError
	}

	return &dto.LoginResponse{Token: token}, "", http.StatusOK
}

func (us *UserService) GetUserFromContext(c *gin.Context) (*models.User, int) {
	cUser, exist := c.Get("user")
	if !exist {
		return nil, http.StatusUnauthorized
	}
	user, ok := cUser.(*models.User)
	if !ok {
		return nil, http.StatusUnauthorized
	}

	return user, http.StatusOK
}

func (us *UserService) getJwtToken(user *models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": user.ID,
		"exp":    time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString([]byte(us.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func NewUserService(userRepo repositories.UserRepository, cartRepo repositories.CartRepository, jwt string) *UserService {
	return &UserService{
		userRepository: userRepo,
		cartRepository: cartRepo,
		jwtSecret:      jwt,
	}
}
