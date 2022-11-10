package postgresql

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/s1ovac/order-subscribe/internal/pkg/utils"
	"github.com/s1ovac/order-subscribe/internal/store/config"
)

type CLient interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
}

func NewClient(ctx context.Context, maxAttempts int, sc *config.StorageConfig) (conn *pgx.Conn, err error) {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", sc.Username, sc.Password, sc.Host, sc.Port, sc.DataBase)
	err = utils.DoWithTries(func() error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		conn, err = pgx.Connect(ctx, dsn)
		if err != nil {
			return err
		}
		return nil
	}, maxAttempts, 5*time.Second)

	if err != nil {
		log.Fatal("error occured connection to postgresql")
	}

	return conn, nil
}
