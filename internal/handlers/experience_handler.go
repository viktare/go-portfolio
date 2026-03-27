package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/viktare/go-portfolio/internal/models"
	"github.com/viktare/go-portfolio/internal/repository"
)

func CreateExperience(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input models.CreateExperience
		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		experience, err := repository.CreateExperience(pool, input)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusCreated, experience)
	}
}

func UpdateExperience(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("experienceID")

		var input models.UpdateExperience
		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		experience, err := repository.UpdateExperience(pool, id, input)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, experience)
	}
}

func GetExperiences(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		experiences, err := repository.GetExperiences(pool)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, experiences)
	}
}

func GetExperience(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("experienceID")

		experience, err := repository.GetExperience(pool, id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, experience)
	}
}

func DeleteExperience(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("experienceID")

		err := repository.DeleteExperience(pool, id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{})
	}
}