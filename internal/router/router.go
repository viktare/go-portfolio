// internal/router/router.go
package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/viktare/go-portfolio/internal/handlers"
)

func Setup(r *gin.Engine, pool *pgxpool.Pool) {
    // Public — no auth
    // public := r.Group("/public")
    // {
    //     public.GET("/portfolio", handlers.GetPortfolio(pool))
    //     public.GET("/resume", handlers.GetResume(pool))
    // }

    // Admin — protected
    admin := r.Group("/admin")
    // admin.Use(AuthMiddleware())
    {
        // User

        // Companies
        admin.GET("/companies", handlers.GetCompanies(pool))
        admin.POST("/companies", handlers.CreateCompany(pool))
		admin.GET("/companies/:companyID", handlers.GetCompany(pool))
        admin.PUT("/companies/:companyID", handlers.UpdateCompany(pool))
        admin.DELETE("/companies/:companyID", handlers.DeleteCompany(pool))

        // Technologies
        admin.GET("/technology-fields", handlers.GetTechnologyFields(pool))
        
        admin.GET("/technologies", handlers.GetTechnologies(pool))
        admin.POST("/technologies", handlers.CreateTechnology(pool))
		admin.GET("/technologies/:technologyID", handlers.GetTechnology(pool))
        admin.PUT("/technologies/:technologyID", handlers.UpdateTechnology(pool))
        admin.DELETE("/technologies/:technologyID", handlers.DeleteTechnology(pool))
	}
}