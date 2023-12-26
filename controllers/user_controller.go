package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zcdanny/task-5-pbi-btpns-Ramanda-Danny/helpers"
	"github.com/zcdanny/task-5-pbi-btpns-Ramanda-Danny/models"
	"gorm.io/gorm"
)

type UserController struct {
	DB *gorm.DB
}

func (uc *UserController) RegisterUser(ctx *gin.Context) {
	var newUser models.User
	if err := ctx.ShouldBindJSON(&newUser); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := helpers.ValidateStruct(newUser); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingUser models.User
	if err := uc.DB.Where("email = ?", newUser.Email).First(&existingUser).Error; err == nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
		return
	}

	hashedPassword, err := helpers.HashPassword(newUser.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	newUser.Password = hashedPassword

	uc.DB.Create(&newUser)

	token, err := helpers.GenerateToken(newUser.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"token": token})
}

func (uc *UserController) LoginUser(ctx *gin.Context) {
	var loginUser models.User
	if err := ctx.ShouldBindJSON(&loginUser); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := helpers.ValidateStruct(loginUser); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingUser models.User
	if err := uc.DB.Where("email = ?", loginUser.Email).First(&existingUser).Error; err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	if !helpers.CheckPasswordHash(loginUser.Password, existingUser.Password) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	token, err := helpers.GenerateToken(existingUser.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

func (uc *UserController) UpdateUser(ctx *gin.Context) {
	userID := ctx.Param("userId")

	var existingUser models.User
	if err := uc.DB.Where("id = ?", userID).First(&existingUser).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var updatedUser models.User
	if err := ctx.ShouldBindJSON(&updatedUser); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := helpers.ValidateStruct(updatedUser); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existingUser.Username = updatedUser.Username
	existingUser.Email = updatedUser.Email
	existingUser.Password = updatedUser.Password

	uc.DB.Save(&existingUser)

	ctx.JSON(http.StatusOK, existingUser)
}

func (uc *UserController) DeleteUser(ctx *gin.Context) {
	userID := ctx.Param("userId")

	var existingUser models.User
	if err := uc.DB.Where("id = ?", userID).First(&existingUser).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	uc.DB.Delete(&existingUser)

	ctx.JSON(http.StatusNoContent, nil)
}