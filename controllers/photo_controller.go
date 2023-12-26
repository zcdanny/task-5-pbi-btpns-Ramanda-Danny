package controllers

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/zcdanny/task-5-pbi-btpns-Ramanda-Danny/helpers"
	"github.com/zcdanny/task-5-pbi-btpns-Ramanda-Danny/models"
	"gorm.io/gorm"
)

type PhotoController struct {
	DB *gorm.DB
}


func (pc *PhotoController) CreatePhoto(ctx *gin.Context) {
	var newPhoto models.Photo

	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get userID from context"})
		return
	}

	userIDString, ok := userID.(string)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to convert userID to string"})
		return
	}

	// Baca title dan caption dari form
	newPhoto.Title = ctx.PostForm("title")
	newPhoto.Caption = ctx.PostForm("caption")

	// Validasi title dan caption tidak boleh kosong
	if newPhoto.Title == "" || newPhoto.Caption == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Title and caption are required"})
		return
	}

	// Validasi struktur data
	if err := helpers.ValidateStruct(newPhoto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Setel UserID dari konteks
	newPhoto.UserID = models.UUIDString(userIDString)

	// Ambil file dari form
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "File upload is required"})
		return
	}

	// Tentukan path penyimpanan file
	filePath := filepath.Join("uploads", string(newPhoto.ID)+filepath.Ext(file.Filename))

	// Simpan file ke path yang ditentukan
	if err := ctx.SaveUploadedFile(file, filePath); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// Setel PhotoUrl dengan path file
	newPhoto.PhotoUrl = filePath

	// Buat data foto baru ke dalam database
	if err := pc.DB.Create(&newPhoto).Error; err != nil {
		fmt.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save to database"})
		return
	}

	// Respon dengan data foto yang berhasil dibuat
	ctx.JSON(http.StatusCreated, newPhoto)
}

func (pc *PhotoController) GetPhotos(ctx *gin.Context) {
	var photos []models.Photo

	pc.DB.Find(&photos)

	ctx.JSON(http.StatusOK, photos)
}

func (pc *PhotoController) GetPhoto(ctx *gin.Context) {
	photoID := ctx.Param("photoId")

	var photo models.Photo
	if err := pc.DB.Where("id = ?", photoID).First(&photo).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Photo not found"})
		return
	}

	ctx.JSON(http.StatusOK, photo)
}

func (pc *PhotoController) UpdatePhoto(ctx *gin.Context) {
	photoID := ctx.Param("photoId")

	var existingPhoto models.Photo
	if err := pc.DB.Where("id = ?", photoID).First(&existingPhoto).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Photo not found"})
		return
	}

	var updatedPhoto models.Photo
	if err := ctx.ShouldBindJSON(&updatedPhoto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := helpers.ValidateStruct(updatedPhoto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existingPhoto.Title = updatedPhoto.Title
	existingPhoto.Caption = updatedPhoto.Caption

	pc.DB.Save(&existingPhoto)

	ctx.JSON(http.StatusOK, existingPhoto)
}

func (pc *PhotoController) DeletePhoto(ctx *gin.Context) {
	photoID := ctx.Param("photoId")

	var existingPhoto models.Photo
	if err := pc.DB.Where("id = ?", photoID).First(&existingPhoto).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Photo not found"})
		return
	}

	pc.DB.Delete(&existingPhoto)

	ctx.JSON(http.StatusNoContent, nil)
}