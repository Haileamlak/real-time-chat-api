package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/haileamlak/chat-system/models"
	"github.com/haileamlak/chat-system/usecases"
)

type UserController interface {
	SignUp(c *gin.Context)
	Login(c *gin.Context)
}

type userController struct {
	userUseCase usecases.UserUseCase
}

func NewUserController(userUseCase usecases.UserUseCase) UserController {
	return &userController{
		userUseCase: userUseCase,
	}
}

func (ctrl *userController) SignUp(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := ctrl.userUseCase.Register(c.Request.Context(), user.Username, user.Password); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

func (ctrl *userController) Login(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	
	token, err := ctrl.userUseCase.Login(c.Request.Context(), user.Username, user.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
	})
}
