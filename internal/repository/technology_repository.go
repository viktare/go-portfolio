package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/viktare/go-portfolio/internal/models"
)

func GetTechnologyFields(pool *pgxpool.Pool) ([]models.TechnologyField, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT id, name, code FROM technology_fields`

	fields := make([]models.TechnologyField, 0)
	rows, err := pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var f models.TechnologyField
		if err := rows.Scan(&f.ID, &f.Name, &f.Code); err != nil {
			return nil, err
		}
		fields = append(fields, f)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return fields, nil
}

func CreateTechnology(pool *pgxpool.Pool, input models.CreateTechnology) (*models.Technology, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO technologies (name, code, field_id)
		VALUES ($1, $2, $3)
		RETURNING id, name, code
	`

	var t models.Technology
	err := pool.QueryRow(ctx, query, input.Name, input.Code, input.FieldID).Scan(
		&t.ID, &t.Name, &t.Code,
	)
	if err != nil {
		return nil, err
	}

	field, err := getTechnologyField(pool, input.FieldID)
	if err != nil {
		return nil, err
	}
	t.Field = *field

	return &t, nil
}

func UpdateTechnology(pool *pgxpool.Pool, technologyID string, input models.UpdateTechnology) (*models.Technology, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		UPDATE technologies
		SET name     = $1,
		    code     = $2,
		    field_id = $3
		WHERE id = $4
		RETURNING id, name, code
	`

	var t models.Technology
	err := pool.QueryRow(ctx, query, input.Name, input.FieldID, technologyID).Scan(
		&t.ID, &t.Name, &t.Code,
	)
	if err != nil {
		return nil, err
	}

	field, err := getTechnologyField(pool, input.FieldID)
	if err != nil {
		return nil, err
	}
	t.Field = *field

	return &t, nil
}

func GetTechnologies(pool *pgxpool.Pool) ([]models.Technology, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT t.id, t.name, t.code,
		       f.id, f.name, f.code
		FROM technologies t
		JOIN technology_fields f ON f.id = t.field_id
	`

	technologies := make([]models.Technology, 0)
	rows, err := pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var t models.Technology
		err := rows.Scan(
			&t.ID, &t.Name, &t.Code,
			&t.Field.ID, &t.Field.Name, &t.Field.Code,
		)
		if err != nil {
			return nil, err
		}
		technologies = append(technologies, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return technologies, nil
}

func GetTechnology(pool *pgxpool.Pool, technologyID string) (*models.Technology, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT t.id, t.name, t.code,
		       f.id, f.name, f.code
		FROM technologies t
		JOIN technology_fields f ON f.id = t.field_id
		WHERE t.id = $1
	`

	var t models.Technology
	err := pool.QueryRow(ctx, query, technologyID).Scan(
		&t.ID, &t.Name, &t.Code,
		&t.Field.ID, &t.Field.Name, &t.Field.Code,
	)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func DeleteTechnology(pool *pgxpool.Pool, technologyID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := pool.Exec(ctx, `DELETE FROM technologies WHERE id = $1`, technologyID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("technology not found")
	}

	return nil
}

func getTechnologyField(pool *pgxpool.Pool, fieldID string) (*models.TechnologyField, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var f models.TechnologyField
	err := pool.QueryRow(ctx, `SELECT id, name, code FROM technology_fields WHERE id = $1`, fieldID).
		Scan(&f.ID, &f.Name, &f.Code)
	if err != nil {
		return nil, err
	}

	return &f, nil
}