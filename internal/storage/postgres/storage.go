package postgres

import (
	"context"
	"fmt"

	"github.com/glebateee/auto-inventory/internal/domain/models"
	"github.com/jackc/pgx/v5"
)

type Storage struct {
	conn *pgx.Conn
}

const (
	dbNotExistsCode = "3D000"
)

func newConnString(dbname string, user string, password string, host string, port int, sslmode string) string {
	return fmt.Sprintf("dbname=%s user=%s password=%s host=%s port=%d sslmode=%s", dbname, user, password, host, port, sslmode)
}
func New(
	dbname string,
	user string,
	password string,
	host string,
	port int,
	sslmode string,
) (*Storage, error) {
	connString := newConnString("postgres", user, password, host, port, sslmode)
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		return nil, err
		// var pgErr pgconn.PgError
		// if ok := errors.As(err, &pgErr); ok && pgErr.Code != dbNotExistsCode {
		// 	return nil, err
		// }
	}
	defer conn.Close(context.Background())
	var exists bool
	err = conn.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname=$1)", dbname).Scan(&exists)
	if err != nil {
		return nil, err
	}
	if !exists {
		fmt.Println("creating db")
		_, err = conn.Exec(context.Background(), fmt.Sprintf("CREATE DATABASE %s", dbname))
	}
	conn.Close(context.Background())

	connString = newConnString(dbname, user, password, host, port, sslmode)
	conn2, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		return nil, err
	}
	return &Storage{
		conn: conn2,
	}, nil
}

func (s *Storage) ProductPageSize(ctx context.Context, page int64, size int64) ([]models.Product, int64) {
	return nil, 0
}
