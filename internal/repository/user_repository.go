package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/viktare/go-portfolio/internal/models"
)

const userDateLayout = "2006-01-02"

func GetUser(pool *pgxpool.Pool) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT u.id, u.name, u.surname, u.bio, u.birth_date, u.location,
		       f.id, f.name, f.path, f.mime_type, f.size
		FROM users u
		LEFT JOIN files f ON f.id = u.resume_file_id
		LIMIT 1
	`

	var (
		u      models.User
		fID    *string
		fName  *string
		fPath  *string
		fMime  *string
		fSize  *int64
	)

	err := pool.QueryRow(ctx, query).Scan(
		&u.ID, &u.Name, &u.Surname, &u.Bio, &u.BirthDate, &u.Location,
		&fID, &fName, &fPath, &fMime, &fSize,
	)
	if err != nil {
		return nil, err
	}

	if fID != nil {
		u.ResumeFile = &models.File{
			ID:       *fID,
			Name:     *fName,
			Path:     *fPath,
			MimeType: *fMime,
			Size:     *fSize,
		}
	}

	contacts, err := GetUserContacts(pool, u.ID)
	if err != nil {
		return nil, err
	}
	u.Contacts = contacts

	experiences, err := GetExperiences(pool)
	if err != nil {
		return nil, err
	}
	u.Experiences = experiences

	projects, err := GetProjects(pool)
	if err != nil {
		return nil, err
	}
	u.Projects = projects

	return &u, nil
}

func UpdateUser(pool *pgxpool.Pool, input models.UpdateUser) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var birthDate *time.Time
	if input.BirthDate != "" {
		t, err := time.Parse(userDateLayout, input.BirthDate)
		if err != nil {
			return nil, fmt.Errorf("invalid birthDate: %w", err)
		}
		birthDate = &t
	}

	query := `
		UPDATE users
		SET name       = $1,
		    surname    = $2,
		    bio        = $3,
		    birth_date = $4,
		    location   = $5
		WHERE id = (SELECT id FROM users LIMIT 1)
		RETURNING id
	`

	var userID string
	err := pool.QueryRow(ctx, query,
		input.Name, input.Surname, input.Bio, birthDate, input.Location,
	).Scan(&userID)
	if err != nil {
		return nil, err
	}

	return GetUser(pool)
}