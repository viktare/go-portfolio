// internal/router/router.go
package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/viktare/go-portfolio/internal/handlers"
)

func Setup(r *gin.Engine, pool *pgxpool.Pool) {

	// Public — no auth
	public := r.Group("/public")
	{
		public.GET("/user", handlers.GetUser(pool))
        public.GET("/user/resume/download", handlers.DownloadResumeFile(pool))
	}

	// Admin — protected
	admin := r.Group("/admin")
	// admin.Use(AuthMiddleware())
	{
		// User
		admin.PUT("/user", handlers.UpdateUser(pool))

		// Resume file
		admin.POST("/user/resume", handlers.UploadResumeFile(pool))
		admin.DELETE("/user/resume", handlers.DeleteResumeFile(pool))

		// User contacts
		admin.POST("/user/contacts", handlers.CreateUserContact(pool))
		admin.PUT("/user/contacts/:contactID", handlers.UpdateUserContact(pool))
		admin.DELETE("/user/contacts/:contactID", handlers.DeleteUserContact(pool))

		// Companies
		admin.GET("/companies", handlers.GetCompanies(pool))
		admin.POST("/companies", handlers.CreateCompany(pool))
		admin.GET("/companies/:companyID", handlers.GetCompany(pool))
		admin.PUT("/companies/:companyID", handlers.UpdateCompany(pool))
		admin.DELETE("/companies/:companyID", handlers.DeleteCompany(pool))

		// Technology fields
		admin.GET("/technology-fields", handlers.GetTechnologyFields(pool))

		// Technologies
		admin.GET("/technologies", handlers.GetTechnologies(pool))
		admin.POST("/technologies", handlers.CreateTechnology(pool))
		admin.GET("/technologies/:technologyID", handlers.GetTechnology(pool))
		admin.PUT("/technologies/:technologyID", handlers.UpdateTechnology(pool))
		admin.DELETE("/technologies/:technologyID", handlers.DeleteTechnology(pool))

		// Projects
		admin.GET("/projects", handlers.GetProjects(pool))
		admin.POST("/projects", handlers.CreateProject(pool))
		admin.GET("/projects/:projectID", handlers.GetProject(pool))
		admin.PUT("/projects/:projectID", handlers.UpdateProject(pool))
		admin.DELETE("/projects/:projectID", handlers.DeleteProject(pool))

		// Experiences
		admin.GET("/experiences", handlers.GetExperiences(pool))
		admin.POST("/experiences", handlers.CreateExperience(pool))
		admin.GET("/experiences/:experienceID", handlers.GetExperience(pool))
		admin.PUT("/experiences/:experienceID", handlers.UpdateExperience(pool))
		admin.DELETE("/experiences/:experienceID", handlers.DeleteExperience(pool))
	}
}