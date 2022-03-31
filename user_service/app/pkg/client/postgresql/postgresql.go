package postgresql

import (
	"context"
	"fmt"
	"github.com/hawkkiller/wtc_system/user_service/internal/config"
	repeatable "github.com/hawkkiller/wtc_system/user_service/pkg/utils"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"time"
)

type Client interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
}

func NewClient(ctx context.Context, sc config.StorageConfig) (pool *pgxpool.Pool, err error) {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		sc.Username, sc.Password, sc.Host, sc.Port, sc.Database,
	)
	err = repeatable.DoWithTries(func() error {
		log.Println("Trying to connect to ", dsn)
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		pool, err = pgxpool.Connect(ctx, dsn)
		if err != nil {
			return err
		}
		return nil
	}, sc.MaxAttempts, 5*time.Second)

	if err != nil {
		log.Fatal("Didn't connect to DB")
		return nil, err
	}
	return
}
