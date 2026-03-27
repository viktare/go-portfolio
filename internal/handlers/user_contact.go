package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/viktare/go-portfolio/internal/models"
	"github.com/viktare/go-portfolio/internal/repository"
)

func CreateUserContact(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input models.CreateUserContact
		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		contact, err := repository.CreateUserContact(pool, input)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusCreated, contact)
	}
}

func UpdateUserContact(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("contactID")

		var input models.UpdateUserContact
		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		contact, err := repository.UpdateUserContact(pool, id, input)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, contact)
	}
}

func DeleteUserContact(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("contactID")

		err := repository.DeleteUserContact(pool, id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{})
	}
}