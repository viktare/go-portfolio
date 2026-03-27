package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/viktare/go-portfolio/internal/models"
)

func UploadResumeFile(pool *pgxpool.Pool, input models.File) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	// Delete old file record if exists (actual disk deletion handled in handler)
	_, err = tx.Exec(ctx, `
		DELETE FROM files
		WHERE id = (SELECT resume_file_id FROM users LIMIT 1)
	`)
	if err != nil {
		return nil, err
	}

	var fileID string
	err = tx.QueryRow(ctx, `
		INSERT INTO files (name, path, mime_type, size)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`, input.Name, input.Path, input.MimeType, input.Size).Scan(&fileID)
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(ctx, `
		UPDATE users SET resume_file_id = $1
		WHERE id = (SELECT id FROM users LIMIT 1)
	`, fileID)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return GetUser(pool)
}

func GetResumeFile(pool *pgxpool.Pool) (*models.File, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    query := `
        SELECT f.id, f.name, f.path, f.mime_type, f.size
        FROM files f
        JOIN users u ON u.resume_file_id = f.id
        LIMIT 1
    `

    var f models.File
    err := pool.QueryRow(ctx, query).Scan(&f.ID, &f.Name, &f.Path, &f.MimeType, &f.Size)
    if err != nil {
        return nil, fmt.Errorf("no resume file found")
    }

    return &f, nil
}

func DeleteResumeFile(pool *pgxpool.Pool) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := pool.Begin(ctx)
	if err != nil {
		return "", err
	}
	defer tx.Rollback(ctx)

	// Grab the path so the handler can delete it from disk
	var filePath string
	err = tx.QueryRow(ctx, `
		SELECT f.path FROM files f
		JOIN users u ON u.resume_file_id = f.id
		LIMIT 1
	`).Scan(&filePath)
	if err != nil {
		return "", fmt.Errorf("no resume file found")
	}

	_, err = tx.Exec(ctx, `UPDATE users SET resume_file_id = NULL WHERE id = (SELECT id FROM users LIMIT 1)`)
	if err != nil {
		return "", err
	}

	_, err = tx.Exec(ctx, `DELETE FROM files WHERE path = $1`, filePath)
	if err != nil {
		return "", err
	}

	if err := tx.Commit(ctx); err != nil {
		return "", err
	}

	return filePath, nil
}