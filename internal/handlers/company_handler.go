package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/viktare/go-portfolio/internal/models"
	"github.com/viktare/go-portfolio/internal/repository"
)

func CreateCompany(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input models.CreateCompanyDTO

		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		product, err := repository.CreateCompany(pool, input)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		ctx.JSON(http.StatusCreated, product)
	}
}

func UpdateCompany(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("companyID")

		var input models.UpdateCompanyDTO
		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		product, err := repository.UpdateCompany(pool, id, input)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, product)
	}
}

func GetCompanies(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		products, err := repository.GetCompanies(pool)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		ctx.JSON(http.StatusOK, products)
	}
}

func GetCompany(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("companyID")
		product, err := repository.GetCompany(pool, id)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		ctx.JSON(http.StatusOK, product)
	}
}

func DeleteCompany(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("companyID")
		err := repository.DeleteCompany(pool, id)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		ctx.JSON(http.StatusOK, gin.H{})
	}
}