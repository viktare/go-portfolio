package handlers

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/viktare/go-portfolio/internal/models"
	"github.com/viktare/go-portfolio/internal/repository"
)

const uploadDir = "./uploads/resume"

func UploadResumeFile(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		file, err := ctx.FormFile("file")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
			return
		}

		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not create upload directory"})
			return
		}

		filePath := filepath.Join(uploadDir, file.Filename)
		if err := ctx.SaveUploadedFile(file, filePath); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not save file"})
			return
		}

		input := models.File{
			Name:     file.Filename,
			Path:     filePath,
			MimeType: file.Header.Get("Content-Type"),
			Size:     file.Size,
		}

		user, err := repository.UploadResumeFile(pool, input)
		if err != nil {
			// Clean up saved file if DB insert fails
			os.Remove(filePath)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, user)
	}
}

func DownloadResumeFile(pool *pgxpool.Pool) gin.HandlerFunc {
    return func(ctx *gin.Context) {
        file, err := repository.GetResumeFile(pool)
        if err != nil {
            ctx.JSON(http.StatusNotFound, gin.H{"error": "no resume file found"})
            return
        }

        ctx.FileAttachment(file.Path, file.Name)
    }
}

func DeleteResumeFile(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		filePath, err := repository.DeleteResumeFile(pool)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "file deleted from db but could not remove from disk"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{})
	}
}