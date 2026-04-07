package v1

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/yertaypert/go-assignment7/internal/entity"
	"github.com/yertaypert/go-assignment7/internal/usecase"
	"github.com/yertaypert/go-assignment7/utils"
)

type userRoutes struct {
	t usecase.UserInterface
}

func RegisterUserRoutes(handler *gin.RouterGroup, t usecase.UserInterface) {
	r := &userRoutes{t: t}
	h := handler.Group("/users")
	{
		h.POST("/", r.RegisterUser)
	}
}

func (r *userRoutes) RegisterUser(c *gin.Context) {
	var createUserDTO entity.CreateUserDTO

	if err := json.NewDecoder(c.Request.Body).Decode(&createUserDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	createUserDTO.Username = strings.TrimSpace(createUserDTO.Username)
	createUserDTO.Email = strings.TrimSpace(createUserDTO.Email)
	createUserDTO.Password = strings.TrimSpace(createUserDTO.Password)
	createUserDTO.Role = strings.TrimSpace(createUserDTO.Role)

	if err := validator.New().Struct(createUserDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := utils.HashPassword(createUserDTO.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}
	user := entity.User{
		Username: createUserDTO.Username,
		Email:    createUserDTO.Email,
		Password: hashedPassword,
		Role:     "user",
	}
	createdUser, sessionID, err := r.t.RegisterUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message":    "User registered successfully. Please check your email for verification code.",
		"session_id": sessionID,
		"user":       createdUser,
	})
}
