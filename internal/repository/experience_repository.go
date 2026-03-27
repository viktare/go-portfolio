package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/viktare/go-portfolio/internal/models"
)

const dateLayout = "2006-01-02"

func CreateExperience(pool *pgxpool.Pool, input models.CreateExperience) (*models.Experience, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	startDate, err := time.Parse(dateLayout, input.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid startDate: %w", err)
	}

	var endDate *time.Time
	if input.EndDate != "" {
		t, err := time.Parse(dateLayout, input.EndDate)
		if err != nil {
			return nil, fmt.Errorf("invalid endDate: %w", err)
		}
		endDate = &t
	}

	tx, err := pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO experiences (role, description, company_id, start_date, end_date, order_index)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, role, description, start_date, end_date, order_index
	`

	var e models.Experience
	err = tx.QueryRow(ctx, query,
		input.Role, input.Description, input.CompanyID, startDate, endDate, input.Order,
	).Scan(&e.ID, &e.Role, &e.Description, &e.StartDate, &e.EndDate, &e.Order)
	if err != nil {
		return nil, err
	}

	if len(input.TechnologyIDs) > 0 {
		if err := insertExperienceTechnologies(ctx, tx, e.ID, input.TechnologyIDs); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	company, err := GetCompany(pool, input.CompanyID)
	if err != nil {
		return nil, err
	}
	e.Company = *company

	technologies, err := getExperienceTechnologies(pool, e.ID)
	if err != nil {
		return nil, err
	}
	e.Technologies = technologies

	return &e, nil
}

func UpdateExperience(pool *pgxpool.Pool, experienceID string, input models.UpdateExperience) (*models.Experience, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	startDate, err := time.Parse(dateLayout, input.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid startDate: %w", err)
	}

	var endDate *time.Time
	if input.EndDate != "" {
		t, err := time.Parse(dateLayout, input.EndDate)
		if err != nil {
			return nil, fmt.Errorf("invalid endDate: %w", err)
		}
		endDate = &t
	}

	tx, err := pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	query := `
		UPDATE experiences
		SET role        = $1,
		    description = $2,
		    company_id  = $3,
		    start_date  = $4,
		    end_date    = $5,
		    order_index = $6
		WHERE id = $7
		RETURNING id, role, description, start_date, end_date, order_index
	`

	var e models.Experience
	err = tx.QueryRow(ctx, query,
		input.Role, input.Description, startDate, endDate, input.Order, experienceID,
	).Scan(&e.ID, &e.Role, &e.Description, &e.StartDate, &e.EndDate, &e.Order)
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(ctx, `DELETE FROM experience_technologies WHERE experience_id = $1`, experienceID)
	if err != nil {
		return nil, err
	}

	if len(input.TechnologyIDs) > 0 {
		if err := insertExperienceTechnologies(ctx, tx, e.ID, input.TechnologyIDs); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	technologies, err := getExperienceTechnologies(pool, e.ID)
	if err != nil {
		return nil, err
	}
	e.Technologies = technologies

	return &e, nil
}

func GetExperiences(pool *pgxpool.Pool) ([]models.Experience, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT e.id, e.role, e.description, e.start_date, e.end_date, e.order_index,
		       c.id, c.name, c.address, c.website, c.logo,
		       t.id, t.name, t.code,
		       f.id, f.name, f.code
		FROM experiences e
		JOIN companies c                   ON c.id = e.company_id
		LEFT JOIN experience_technologies et ON et.experience_id = e.id
		LEFT JOIN technologies t           ON t.id = et.technology_id
		LEFT JOIN technology_fields f      ON f.id = t.field_id
		ORDER BY e.order_index ASC
	`

	rows, err := pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	experienceMap := make(map[string]*models.Experience)
	experienceOrder := make([]string, 0)

	for rows.Next() {
		var (
			e     models.Experience
			tID   *string
			tName *string
			tCode *string
			fID   *string
			fName *string
			fCode *string
		)

		err := rows.Scan(
			&e.ID, &e.Role, &e.Description, &e.StartDate, &e.EndDate, &e.Order,
			&e.Company.ID, &e.Company.Name, &e.Company.Address, &e.Company.Website, &e.Company.Logo,
			&tID, &tName, &tCode,
			&fID, &fName, &fCode,
		)
		if err != nil {
			return nil, err
		}

		if _, exists := experienceMap[e.ID]; !exists {
			e.Technologies = make([]models.Technology, 0)
			experienceMap[e.ID] = &e
			experienceOrder = append(experienceOrder, e.ID)
		}

		if tID != nil {
			experienceMap[e.ID].Technologies = append(experienceMap[e.ID].Technologies, models.Technology{
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

	experiences := make([]models.Experience, 0, len(experienceOrder))
	for _, id := range experienceOrder {
		experiences = append(experiences, *experienceMap[id])
	}

	return experiences, nil
}

func GetExperience(pool *pgxpool.Pool, experienceID string) (*models.Experience, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT e.id, e.role, e.description, e.start_date, e.end_date, e.order_index,
		       c.id, c.name, c.address, c.website, c.logo,
		       t.id, t.name, t.code,
		       f.id, f.name, f.code
		FROM experiences e
		JOIN companies c                     ON c.id = e.company_id
		LEFT JOIN experience_technologies et ON et.experience_id = e.id
		LEFT JOIN technologies t             ON t.id = et.technology_id
		LEFT JOIN technology_fields f        ON f.id = t.field_id
		WHERE e.id = $1
	`

	rows, err := pool.Query(ctx, query, experienceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var e *models.Experience

	for rows.Next() {
		var (
			row   models.Experience
			tID   *string
			tName *string
			tCode *string
			fID   *string
			fName *string
			fCode *string
		)

		err := rows.Scan(
			&row.ID, &row.Role, &row.Description, &row.StartDate, &row.EndDate, &row.Order,
			&row.Company.ID, &row.Company.Name, &row.Company.Address, &row.Company.Website, &row.Company.Logo,
			&tID, &tName, &tCode,
			&fID, &fName, &fCode,
		)
		if err != nil {
			return nil, err
		}

		if e == nil {
			row.Technologies = make([]models.Technology, 0)
			e = &row
		}

		if tID != nil {
			e.Technologies = append(e.Technologies, models.Technology{
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

	if e == nil {
		return nil, fmt.Errorf("experience not found")
	}

	return e, nil
}

func DeleteExperience(pool *pgxpool.Pool, experienceID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := pool.Exec(ctx, `DELETE FROM experiences WHERE id = $1`, experienceID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("experience not found")
	}

	return nil
}

// =============================================================================
// Helpers
// =============================================================================

func insertExperienceTechnologies(ctx context.Context, tx pgx.Tx, experienceID string, technologyIDs []string) error {
	for _, tID := range technologyIDs {
		_, err := tx.Exec(ctx,
			`INSERT INTO experience_technologies (experience_id, technology_id) VALUES ($1, $2)`,
			experienceID, tID,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func getExperienceTechnologies(pool *pgxpool.Pool, experienceID string) ([]models.Technology, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT t.id, t.name, t.code,
		       f.id, f.name, f.code
		FROM experience_technologies et
		JOIN technologies t      ON t.id = et.technology_id
		JOIN technology_fields f ON f.id = t.field_id
		WHERE et.experience_id = $1
	`

	rows, err := pool.Query(ctx, query, experienceID)
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