package db

import (
	"context"
	"fmt"
	"github.com/hawkkiller/wtc_system/user_service/internal/apperror"
	"github.com/hawkkiller/wtc_system/user_service/internal/user"
	"github.com/hawkkiller/wtc_system/user_service/pkg/client/postgresql"
	"github.com/hawkkiller/wtc_system/user_service/pkg/logging"
	"github.com/jackc/pgconn"
	"strings"
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
		select id, username, email, password from "user" where email = $1
		`
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

func (r *db) Create(ctx context.Context, user *user.User) error {
	q := `
		insert into "user" (username, password, email) 
		values ($1, $2, $3) 
		returning id
	`
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
				return apperror.NewAppError(
					"That email is already registered in the system",
					"WTC-000004",
					"enter another email",
				)
			}
			return nil
		}
		return err
	}
	return nil
}

func (r *db) FindOne(ctx context.Context, id string) (user.User, error) {
	q := `
		select id, username, email, password from "user" where id = $1
		`
	var u user.User

	err := r.client.QueryRow(ctx, q, id).Scan(&u.ID, &u.Username, &u.Email, &u.Password)
	if err != nil {
		return u, err
	}

	return u, nil

}

func (r *db) Update(ctx context.Context, user user.User) error {
	//TODO implement me
	panic("implement me")
}

func (r *db) Delete(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}
