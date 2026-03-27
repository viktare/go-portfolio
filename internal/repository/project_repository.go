package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/viktare/go-portfolio/internal/models"
)

func CreateProject(pool *pgxpool.Pool, input models.CreateProject) (*models.Project, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO projects (name, code, status, description, link)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, name, code, status, description, link
	`

	var p models.Project
	err = tx.QueryRow(ctx, query,
		input.Name, input.Code, input.Status, input.Description, input.Link,
	).Scan(&p.ID, &p.Name, &p.Code, &p.Status, &p.Description, &p.Link)
	if err != nil {
		return nil, err
	}

	if len(input.TechnologyIDs) > 0 {
		if err := insertProjectTechnologies(ctx, tx, p.ID, input.TechnologyIDs); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	technologies, err := getProjectTechnologies(pool, p.ID)
	if err != nil {
		return nil, err
	}
	p.Technologies = technologies

	return &p, nil
}

func UpdateProject(pool *pgxpool.Pool, projectID string, input models.UpdateProject) (*models.Project, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	query := `
		UPDATE projects
		SET name        = $1,
		    code        = $2,
		    status      = $3,
		    description = $4,
		    link        = $5
		WHERE id = $6
		RETURNING id, name, code, status, description, link
	`

	var p models.Project
	err = tx.QueryRow(ctx, query,
		input.Name, input.Status, input.Description, input.Link, projectID,
	).Scan(&p.ID, &p.Name, &p.Code, &p.Status, &p.Description, &p.Link)
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(ctx, `DELETE FROM project_technologies WHERE project_id = $1`, projectID)
	if err != nil {
		return nil, err
	}

	if len(input.TechnologyIDs) > 0 {
		if err := insertProjectTechnologies(ctx, tx, p.ID, input.TechnologyIDs); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	technologies, err := getProjectTechnologies(pool, p.ID)
	if err != nil {
		return nil, err
	}
	p.Technologies = technologies

	return &p, nil
}

func GetProjects(pool *pgxpool.Pool) ([]models.Project, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT p.id, p.name, p.code, p.status, p.description, p.link,
		       t.id, t.name, t.code,
		       f.id, f.name, f.code
		FROM projects p
		LEFT JOIN project_technologies pt ON pt.project_id = p.id
		LEFT JOIN technologies t          ON t.id = pt.technology_id
		LEFT JOIN technology_fields f     ON f.id = t.field_id
		ORDER BY p.id
	`

	rows, err := pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	projectMap := make(map[string]*models.Project)
	projectOrder := make([]string, 0)

	for rows.Next() {
		var (
			p     models.Project
			tID   *string
			tName *string
			tCode *string
			fID   *string
			fName *string
			fCode *string
		)

		err := rows.Scan(
			&p.ID, &p.Name, &p.Code, &p.Status, &p.Description, &p.Link,
			&tID, &tName, &tCode,
			&fID, &fName, &fCode,
		)
		if err != nil {
			return nil, err
		}

		if _, exists := projectMap[p.ID]; !exists {
			p.Technologies = make([]models.Technology, 0)
			projectMap[p.ID] = &p
			projectOrder = append(projectOrder, p.ID)
		}

		if tID != nil {
			projectMap[p.ID].Technologies = append(projectMap[p.ID].Technologies, models.Technology{
				ID:   *tID,
				Name: *tName,
				Code: *tCode,
				Field: models.TechnologyField{
					ID:   *fID,
					Name: *fName,
					Code: *fCode,
				},
			})
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	projects := make([]models.Project, 0, len(projectOrder))
	for _, id := range projectOrder {
		projects = append(projects, *projectMap[id])
	}

	return projects, nil
}

func GetProject(pool *pgxpool.Pool, projectID string) (*models.Project, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT p.id, p.name, p.code, p.status, p.description, p.link,
		       t.id, t.name, t.code,
		       f.id, f.name, f.code
		FROM projects p
		LEFT JOIN project_technologies pt ON pt.project_id = p.id
		LEFT JOIN technologies t          ON t.id = pt.technology_id
		LEFT JOIN technology_fields f     ON f.id = t.field_id
		WHERE p.id = $1
	`

	rows, err := pool.Query(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var p *models.Project

	for rows.Next() {
		var (
			row   models.Project
			tID   *string
			tName *string
			tCode *string
			fID   *string
			fName *string
			fCode *string
		)

		err := rows.Scan(
			&row.ID, &row.Name, &row.Code, &row.Status, &row.Description, &row.Link,
			&tID, &tName, &tCode,
			&fID, &fName, &fCode,
		)
		if err != nil {
			return nil, err
		}

		if p == nil {
			row.Technologies = make([]models.Technology, 0)
			p = &row
		}

		if tID != nil {
			p.Technologies = append(p.Technologies, models.Technology{
				ID:   *tID,
				Name: *tName,
				Code: *tCode,
				Field: models.TechnologyField{
					ID:   *fID,
					Name: *fName,
					Code: *fCode,
				},
			})
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if p == nil {
		return nil, fmt.Errorf("project not found")
	}

	return p, nil
}

func DeleteProject(pool *pgxpool.Pool, projectID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := pool.Exec(ctx, `DELETE FROM projects WHERE id = $1`, projectID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("project not found")
	}

	return nil
}

func insertProjectTechnologies(ctx context.Context, tx pgx.Tx, projectID string, technologyIDs []string) error {
	for _, tID := range technologyIDs {
		_, err := tx.Exec(ctx,
			`INSERT INTO project_technologies (project_id, technology_id) VALUES ($1, $2)`,
			projectID, tID,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func getProjectTechnologies(pool *pgxpool.Pool, projectID string) ([]models.Technology, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT t.id, t.name, t.code,
		       f.id, f.name, f.code
		FROM project_technologies pt
		JOIN technologies t      ON t.id = pt.technology_id
		JOIN technology_fields f ON f.id = t.field_id
		WHERE pt.project_id = $1
	`

	rows, err := pool.Query(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	technologies := make([]models.Technology, 0)
	for rows.Next() {
		var t models.Technology
		if err := rows.Scan(&t.ID, &t.Name, &t.Code, &t.Field.ID, &t.Field.Name, &t.Field.Code); err != nil {
			return nil, err
		}
		technologies = append(technologies, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return technologies, nil
}