package db

import (
	"context"
	"fmt"
	"github.com/WTC-SYSTEM/wtc_system/libs/logging"
	"github.com/WTC-SYSTEM/wtc_system/user_service/internal/apperror"
	"github.com/WTC-SYSTEM/wtc_system/user_service/internal/user"
	"github.com/WTC-SYSTEM/wtc_system/user_service/pkg/client/postgresql"
	"github.com/jackc/pgconn"
	"strings"
	"time"
)

type db struct {
	client postgresql.Client
	logger logging.Logger
}

func NewStorage(c postgresql.Client, l logging.Logger) user.Storage {
	return &db{
		client: c,
		logger: l,
	}
}

func formatQuery(q string) string {
	return strings.ReplaceAll(strings.ReplaceAll(q, "\t", ""), "\n", "")
}

func (r *db) FindByEmail(ctx context.Context, email string) (user.User, error) {
	q := `
		SELECT "id", "username", "email", "password" 
		FROM "users" WHERE "email" = $1
		`
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	var u user.User
	r.logger.Trace(fmt.Sprintf("SQL query: %s", formatQuery(q)))

	if err := r.client.QueryRow(ctx, q, email).Scan(&u.ID, &u.Username, &u.Email, &u.Password); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			r.logger.Error(
				"SQL Error: ", pgErr.Message,
				", Detail: ", pgErr.Detail,
				", Where: ", pgErr.Where,
				", Code: ", pgErr.Code,
				", SQLState: ", pgErr.SQLState(),
			)

		}
		return user.User{}, err
	}
	return u, nil
}

func (r *db) Create(ctx context.Context, user *user.User) (string, error) {
	q := `
		INSERT INTO "users" ("username", "password", "email") 
		VALUES ($1, $2, $3) 
		RETURNING id
	`
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	r.logger.Trace(fmt.Sprintf("SQL query: %s", formatQuery(q)))
	if err := r.client.QueryRow(ctx, q, user.Username, user.Password, user.Email).Scan(&user.ID); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			r.logger.Error(
				"SQL Error: ", pgErr.Message,
				", Detail: ", pgErr.Detail,
				", Where: ", pgErr.Where,
				", Code: ", pgErr.Code,
				", SQLState: ", pgErr.SQLState(),
			)
			if strings.Contains(pgErr.Message, "duplicate key value violates unique") {
				return "", apperror.NewAppError(
					"That email is already registered in the system",
					"WTC-000004",
					"enter another email",
				)
			}
			return "", err
		}
		return "", err
	}
	return user.ID, nil
}

func (r *db) FindOne(ctx context.Context, id string) (user.User, error) {
	q := `
		SELECT "id", "username", "email", "password" 
		FROM "users" where id = $1
		`
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	var u user.User

	err := r.client.QueryRow(ctx, q, id).Scan(&u.ID, &u.Username, &u.Email, &u.Password)
	if err != nil {
		return u, err
	}

	return u, nil

}

func (r *db) Update(ctx context.Context, user user.User) error {
	q := `
		UPDATE users
		SET username = $1, email = $2, password = $3
		WHERE id = $4;
		`
	r.logger.Trace(q)
	r.logger.Trace(user)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := r.client.QueryRow(ctx, q, user.Username, user.Email, user.Password, user.ID).Scan(); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			r.logger.Error(
				"SQL Error: ", pgErr.Message,
				", Detail: ", pgErr.Detail,
				", Where: ", pgErr.Where,
				", Code: ", pgErr.Code,
				", SQLState: ", pgErr.SQLState(),
			)
			if strings.Contains(pgErr.Message, "duplicate key value violates unique") {
				return apperror.NewAppError(
					"That email is already registered in the system",
					"WTC-000004",
					"enter another email",
				)
			}
			return err
		}
		return err
	}

	return nil
}

func (r *db) Delete(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}
