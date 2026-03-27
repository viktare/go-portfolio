package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/viktare/go-portfolio/internal/models"
)

func CreateCompany(pool *pgxpool.Pool, input models.CreateCompanyDTO) (*models.Company, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	query := `
		INSERT INTO companies (name, address, website, logo)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, address, website, logo
	`
	var company models.Company

	err := pool.QueryRow(ctx, query, input.Name, input.Address, input.Website, input.Logo).Scan(
		&company.ID,
		&company.Name,
		&company.Address,
		&company.Website,
		&company.Logo,
	)

	if err != nil {
		return nil, err
	}

	return &company, nil
}

func UpdateCompany(pool *pgxpool.Pool, input models.UpdateCompanyDTO) (*models.Company, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	query := `
		UPDATE companies
		SET name = $1,
			address = $2,
			website = $3,
			logo = $4 
		RETURNING id, name, address, website, logo
	`
	var company models.Company

	err := pool.QueryRow(ctx, query, input.Name, input.Address, input.Website, input.Logo).Scan(
		&company.ID,
		&company.Name,
		&company.Address,
		&company.Website,
		&company.Logo,
	)

	if err != nil {
		return nil, err
	}

	return &company, nil
}

func GetCompanies(pool *pgxpool.Pool) ([]models.Company, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	query := `
		SELECT id, name, address, website, logo
		FROM companies
	`

	companies := make([]models.Company, 0)
	rows, err := pool.Query(ctx, query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var c models.Company
		err := rows.Scan(&c.ID, &c.Name, &c.Address, &c.Website, &c.Logo)
		if err != nil {
			return nil, err
		}
		companies = append(companies, c)
	}

	// check for errors that occurred during iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return companies, nil
}

func GetCompany(pool *pgxpool.Pool, companyID string) (*models.Company, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    query := `
        SELECT id, name, address, website, logo
        FROM companies
        WHERE id = $1
    `

    var c models.Company
    err := pool.QueryRow(ctx, query, companyID).Scan(
        &c.ID, &c.Name, &c.Address, &c.Website, &c.Logo,
    )
    if err != nil {
        return nil, err
    }

    return &c, nil
}

func DeleteCompany(pool *pgxpool.Pool, companyID string) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    query := `DELETE FROM companies WHERE id = $1`

    result, err := pool.Exec(ctx, query, companyID)
    if err != nil {
        return err
    }

    if result.RowsAffected() == 0 {
        return fmt.Errorf("Company not found")
    }

    return nil
}