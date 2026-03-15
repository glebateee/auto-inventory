package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/glebateee/auto-inventory/internal/domain/models"
	"github.com/glebateee/auto-inventory/internal/storage"
	sqlc "github.com/glebateee/auto-inventory/internal/storage/postgres/sqlc/gen"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type Storage struct {
	querier sqlc.Querier
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
		querier: sqlc.New(conn2),
	}, nil
}

func (s *Storage) ProductPageSizeCategory(
	ctx context.Context,
	offset int64,
	limit int64,
	categoryID int64,
) ([]models.Product, int64, error) {
	pgCategoryId := pgtype.Int4{Int32: int32(categoryID), Valid: true}
	sqlcTotal, err := s.querier.ProductTotalCategory(ctx, pgCategoryId)
	fmt.Println(sqlcTotal, categoryID)
	if err != nil {
		return nil, 0, err
	}
	sqlcProducts, err := s.querier.ProductPageSizeCategory(ctx, sqlc.ProductPageSizeCategoryParams{
		Offset: int32(offset),
		Limit:  int32(limit),
		ID:     int32(categoryID),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sqlcTotal, storage.ErrNoRows
		}
		return nil, sqlcTotal, err
	}
	products := FromSqlcProductListCat(sqlcProducts)
	return products, sqlcTotal, nil
}

func (s *Storage) ProductPageSize(ctx context.Context, page int64, size int64) ([]models.Product, int64, error) {
	sqlcTotal, err := s.querier.ProductTotal(ctx)
	if err != nil {
		return nil, 0, err
	}
	sqlcProducts, err := s.querier.ProductPageSize(ctx, sqlc.ProductPageSizeParams{Offset: int32((page - 1) * size), Limit: int32(size)})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sqlcTotal, storage.ErrNoRows
		}
		return nil, sqlcTotal, err
	}
	products := FromSqlcProductList(sqlcProducts)
	return products, sqlcTotal, nil
}

func newConnString(dbname string, user string, password string, host string, port int, sslmode string) string {
	return fmt.Sprintf("dbname=%s user=%s password=%s host=%s port=%d sslmode=%s", dbname, user, password, host, port, sslmode)
}
