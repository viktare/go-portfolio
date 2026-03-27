package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/viktare/go-portfolio/internal/models"
)

func GetUserContacts(pool *pgxpool.Pool, userID string) ([]models.UserContact, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, user_id, name, code, link
		FROM user_contacts
		WHERE user_id = $1
	`

	rows, err := pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	contacts := make([]models.UserContact, 0)
	for rows.Next() {
		var c models.UserContact
		if err := rows.Scan(&c.ID, &c.UserID, &c.Name, &c.Code, &c.Link); err != nil {
			return nil, err
		}
		contacts = append(contacts, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return contacts, nil
}

func CreateUserContact(pool *pgxpool.Pool, input models.CreateUserContact) (*models.UserContact, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO user_contacts (user_id, name, code, link)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, name, code, link
	`

	var c models.UserContact
	err := pool.QueryRow(ctx, query, input.UserID, input.Name, input.Code, input.Link).Scan(
		&c.ID, &c.UserID, &c.Name, &c.Code, &c.Link,
	)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func UpdateUserContact(pool *pgxpool.Pool, contactID string, input models.UpdateUserContact) (*models.UserContact, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		UPDATE user_contacts
		SET name = $1,
		    code = $2,
		    link = $3
		WHERE id = $4
		RETURNING id, user_id, name, code, link
	`

	var c models.UserContact
	err := pool.QueryRow(ctx, query, input.Name, input.Link, contactID).Scan(
		&c.ID, &c.UserID, &c.Name, &c.Code, &c.Link,
	)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func DeleteUserContact(pool *pgxpool.Pool, contactID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := pool.Exec(ctx, `DELETE FROM user_contacts WHERE id = $1`, contactID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("contact not found")
	}

	return nil
}