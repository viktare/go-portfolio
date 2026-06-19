// internal/handlers/public.go
package handlers

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/viktare/go-portfolio/internal/repository"
)

func GetPortfolio(pool *pgxpool.Pool) gin.HandlerFunc {
    // Parse templates once at startup, not on every request
    tmpl := template.Must(template.ParseGlob("templates/**/*.html"))

    return func(c *gin.Context) {
        user, err := repository.GetUser(pool)
        if err != nil {
            c.String(http.StatusInternalServerError, "failed to load portfolio")
            return
        }

        c.Header("Content-Type", "text/html; charset=utf-8")
        tmpl.ExecuteTemplate(c.Writer, "base", user)
    }
}