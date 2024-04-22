package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type UserModel struct {
	DB *sql.DB
}

var (
	ErrDuplicateEmail = errors.New("duplicate email")
)

func (m UserModel) Insert(user *User) error {
	query := `INSERT INTO user_info (fname, sname, email, password_hash, user_role, activated)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id, created_at, updated_at, version`
	args := []interface{}{user.FName, user.SName, user.Email, user.Password.hash, user.UserRole, user.Activated}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}
	return nil
}

func (m UserModel) GetByEmail(email string) (*User, error) {
	query := `SELECT id, created_at, updated_at, fname, sname, email, password_hash, user_role, activated, version
    FROM user_info
    WHERE email = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	var user User

	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.FName,
		&user.SName,
		&user.Email,
		&user.Password.hash,
		&user.UserRole,
		&user.Activated,
		&user.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}

func (m UserModel) Update(user *User) error {
	query := `
UPDATE user_info
SET fname = $1, sname = $2, email = $3, password_hash = $4, activated=$5, version = version + 1
WHERE id = $6 AND version = $7
RETURNING version`
	args := []any{
		user.FName,
		user.SName,
		user.Email,
		user.Password.hash,
		user.Activated,
		user.ID,
		user.Version,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		case errors.Is(err, sql.ErrNoRows):
			return nil
		default:
			return err
		}
	}
	return nil
}

func (m UserModel) GetUser(id int64) (*User, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := `SELECT id, created_at, fname, sname, email, user_role, activated, version FROM user_info WHERE id = $1`
	var user User

	err := m.DB.QueryRow(query, id).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.FName,
		&user.SName,
		&user.Email,
		&user.UserRole,
		&user.Activated,
		&user.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}

func (m UserModel) DeleteUser(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `DELETE FROM user_info WHERE id = $1`
	result, err := m.DB.Exec(query, id)

	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil

}
