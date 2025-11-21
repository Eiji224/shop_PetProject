package services

import (
	"context"
	"errors"
	"net/http"
	"shop/internal/dto"
	"shop/internal/models"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

type MockUserRepository struct {
	mock.Mock
}

func (repo *MockUserRepository) Insert(user *models.User, ctx context.Context) error {
	return repo.Called(user, ctx).Error(0)
}

func (repo *MockUserRepository) FindByID(id uint, ctx context.Context) (*models.User, error) {
	args := repo.Called(id, ctx)
	return args.Get(0).(*models.User), args.Error(1)
}

func (repo *MockUserRepository) FindByEmail(email string, ctx context.Context) (*models.User, error) {
	args := repo.Called(email, ctx)
	return args.Get(0).(*models.User), args.Error(1)
}

func newTestUserService() (*UserService, *MockUserRepository, *MockCartRepository, string) {
	jwtSecret := "testing"
	userRepo := new(MockUserRepository)
	cartRepo := new(MockCartRepository)
	service := NewUserService(userRepo, cartRepo, jwtSecret)
	return service, userRepo, cartRepo, jwtSecret
}

func getTestUser() *models.User {
	testPassword := "andrey"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(testPassword), bcrypt.DefaultCost)

	return &models.User{
		Username: "Andrey",
		Email:    "andrey@example.com",
		Password: string(hashedPassword),
		Type:     "customer",
	}
}

func TestCreate(t *testing.T) {
	userService, userRepo, cartRepo, jwtSecret := newTestUserService()

	req := dto.RegisterRequest{
		Username: "Andrey",
		Email:    "andrey@example.com",
		Password: "Andrey",
		Type:     "customer",
	}

	user := getTestUser()

	userRepo.On("Insert",
		mock.MatchedBy(func(u *models.User) bool {
			return u.Email == user.Email &&
				u.Username == user.Username &&
				u.Type == user.Type
			// Пароль и ID игнорируем — они генерируются/хэшируются внутри
		}),
		mock.Anything,
	).Return(nil)
	cartRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

	resp, errMsg, status := userService.Register(req)

	assert.Equal(t, http.StatusCreated, status)
	assert.Equal(t, "", errMsg)
	assert.Equal(t, req.Email, resp.Email)
	assert.Equal(t, user.Username, resp.Username)
	assert.Equal(t, true, checkToken(resp.Token, jwtSecret))

	userRepo.AssertExpectations(t)
	cartRepo.AssertExpectations(t)
}

func TestSuccessfulLogin(t *testing.T) {
	userService, userRepo, _, jwtSecret := newTestUserService()

	loginReq := dto.LoginRequest{
		Email:    "andrey@example.com",
		Password: "andrey",
	}
	user := getTestUser()

	userRepo.On("FindByEmail", "andrey@example.com", mock.Anything).Return(user, nil)

	resp, errMsg, status := userService.Login(loginReq)

	assert.Equal(t, http.StatusOK, status)
	assert.Equal(t, "", errMsg)
	assert.Equal(t, true, checkToken(resp.Token, jwtSecret))

	userRepo.AssertExpectations(t)
}

func TestFailedLogin(t *testing.T) {
	userService, userRepo, _, _ := newTestUserService()

	loginReqs := []dto.LoginRequest{
		{
			Email:    "somebody@example.com",
			Password: "123",
		},
		{
			Email:    "andrey@example.com",
			Password: "123",
		},
	}

	user := getTestUser()

	userRepo.On("FindByEmail", "somebody@example.com", mock.Anything).Return((*models.User)(nil), errors.New("user not found"))
	userRepo.On("FindByEmail", "andrey@example.com", mock.Anything).Return(user, nil)

	for _, req := range loginReqs {
		resp, errMsg, status := userService.Login(req)
		assert.Equal(t, http.StatusUnauthorized, status)
		assert.Equal(t, "Invalid username or password", errMsg)
		assert.Nil(t, resp)
	}

	userRepo.AssertExpectations(t)
}

func checkToken(tokenString string, jwtSecret string) bool {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}

		return []byte(jwtSecret), nil
	})
	if err != nil || !token.Valid {
		return false
	}

	_, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false
	}

	return true
}
