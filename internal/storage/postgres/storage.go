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
	conn    *pgx.Conn
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

// func (s *Storage) UpdateProduct(ctx context.Context, sku string, fields *models.UpdateProductFields, mask *fieldmaskpb.FieldMask) (*models.Product, error) {
// 	tx, err := s.conn.Begin(ctx)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer tx.Rollback(ctx)
// 	setClauses := []string{}
// 	args := pgx.NamedArgs{"sku": sku}
// 	for _, path := range mask.GetPaths() {
// 		switch path {
// 		case "name":
// 			if fields.Name != nil {
// 				setClauses = append(setClauses, "name = @name")
// 				args["name"] = *fields.Name
// 			}
// 		case "description":
// 			if fields.Description != nil {
// 				setClauses = append(setClauses, "description = @description")
// 				args["description"] = *fields.Description
// 			}
// 		case "category":
// 			if fields.Category != nil {
// 				var categoryID int32
// 				err = tx.QueryRow(ctx, "SELECT id FROM categories WHERE name = $1", *fields.Category).Scan(&categoryID)
// 				if err != nil {
// 					return nil, fmt.Errorf("category not found: %w", err)
// 				}
// 				setClauses = append(setClauses, "category_id = @category_id")
// 				args["category_id"] = categoryID
// 			}
// 		case "manufacturer":
// 			if fields.Manufacturer != nil {
// 				var manufacturerID int32
// 				err = tx.QueryRow(ctx, "SELECT id FROM manufacturers WHERE name = $1", *fields.Manufacturer).Scan(&manufacturerID)
// 				if err != nil {
// 					return nil, fmt.Errorf("manufacturer not found: %w", err)
// 				}
// 				setClauses = append(setClauses, "manufacturer_id = @manufacturer_id")
// 				args["manufacturer_id"] = manufacturerID
// 			}
// 		case "weight":
// 			if fields.Weight != nil {
// 				setClauses = append(setClauses, "weight = @weight")
// 				args["weight"] = *fields.Weight
// 			}
// 		case "price":
// 			if fields.Price != nil {
// 				setClauses = append(setClauses, "price = @price")
// 				args["price"] = *fields.Price
// 			}
// 		case "base_price":
// 			if fields.BasePrice != nil {
// 				setClauses = append(setClauses, "baseprice = @baseprice")
// 				args["baseprice"] = *fields.BasePrice
// 			}
// 		case "issue_year":
// 			if fields.IssueYear != nil {
// 				setClauses = append(setClauses, "issueyear = @issueyear")
// 				args["issueyear"] = *fields.IssueYear
// 			}
// 		}
// 	}
// 	// if len(setClauses) == 0 {
// 	// 	// Nothing to update – return current product
// 	// 	return s.getProductBySku(ctx, sku) // you'll need to implement this helper
// 	// }
// 	// Always update updated_at (or rely on trigger)
// 	setClauses = append(setClauses, "updated_at = NOW()")
// 	query := fmt.Sprintf(`
//         UPDATE products
//         SET %s
//         WHERE sku = @sku
//         RETURNING *
//     `, strings.Join(setClauses, ", "))
// 	row := tx.QueryRow(ctx, query, args)
// 	var product models.Product
// 	if err := scanProduct(row, &product); err != nil {
// 		return nil, err
// 	}
//		if err := tx.Commit(ctx); err != nil {
//			return nil, err
//		}
//		return &product, nil
//	}
//
//	func scanProduct(row pgx.Row, p *models.Product) error {
//		var (
//			categoryName, manufacturerName string
//			description                    sql.NullString
//			unit                           sql.NullString
//			createdAt, updatedAt           time.Time
//		)
//		err := row.Scan(
//			&p.Id,
//			&p.Sku,
//			&p.Name,
//			&description,
//			&categoryName,
//			&manufacturerName,
//			&p.Weight,
//			&unit,
//			&p.Price,
//			&p.BasePrice,
//			&p.IssueYear,
//			&createdAt,
//			&updatedAt,
//			// add any other columns from your products table (including created_at, updated_at)
//		)
//		if err != nil {
//			return err
//		}
//		p.Description = description.String
//		p.Category = categoryName
//		p.Manufacturer = manufacturerName
//		p.CreatedAt = createdAt
//		p.UpdatedAt = updatedAt
//		// unit is not in your domain model? If needed, add it.
//		return nil
//	}

func (s *Storage) DeleteProductSku(ctx context.Context, sku string) error {
	const op = "postgres.DeleteProductSku"
	affected, err := s.querier.DeleteProductBySku(ctx, sku)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if affected != 1 {
		if affected == 0 {
			return fmt.Errorf("%s: %w", op, storage.ErrNoRows)
		}
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Storage) Products(ctx context.Context) ([]models.Product, error) {
	const op = "postgres.Products"
	sqlcProducts, err := s.querier.Products(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []models.Product{}, fmt.Errorf("%s: %w", op, storage.ErrNoRows)
		}
		return []models.Product{}, fmt.Errorf("%s: %w", op, err)
	}
	return FromSqlcProducts(sqlcProducts), nil
}

func (s *Storage) ProductPageSizeCategory(
	ctx context.Context,
	offset int64,
	limit int64,
	categoryID int64,
) ([]models.Product, int64, error) {
	const op = "postgres.ProductPageSizeCategory"
	pgCategoryId := pgtype.Int4{Int32: int32(categoryID), Valid: true}
	sqlcTotal, err := s.querier.ProductTotalCategory(ctx, pgCategoryId)
	fmt.Println(sqlcTotal, categoryID)
	if err != nil {
		return []models.Product{}, 0, fmt.Errorf("%s: %w", op, err)
	}
	sqlcProducts, err := s.querier.ProductPageSizeCategory(ctx, sqlc.ProductPageSizeCategoryParams{
		Offset: int32(offset),
		Limit:  int32(limit),
		ID:     int32(categoryID),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []models.Product{}, sqlcTotal, fmt.Errorf("%s: %w", op, storage.ErrNoRows)
		}
		return []models.Product{}, sqlcTotal, fmt.Errorf("%s: %w", op, err)
	}
	products := FromSqlcProductListCat(sqlcProducts)
	return products, sqlcTotal, nil
}

func (s *Storage) ProductPageSize(ctx context.Context, page int64, size int64) ([]models.Product, int64, error) {
	const op = "postgres.ProductPageSize"
	sqlcTotal, err := s.querier.ProductTotal(ctx)
	if err != nil {
		return []models.Product{}, 0, fmt.Errorf("%s: %w", op, err)
	}
	sqlcProducts, err := s.querier.ProductPageSize(ctx, sqlc.ProductPageSizeParams{Offset: int32((page - 1) * size), Limit: int32(size)})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []models.Product{}, sqlcTotal, fmt.Errorf("%s: %w", op, storage.ErrNoRows)
		}
		return []models.Product{}, sqlcTotal, fmt.Errorf("%s: %w", op, err)
	}
	products := FromSqlcProductList(sqlcProducts)
	return products, sqlcTotal, nil
}

func newConnString(dbname string, user string, password string, host string, port int, sslmode string) string {
	return fmt.Sprintf("dbname=%s user=%s password=%s host=%s port=%d sslmode=%s", dbname, user, password, host, port, sslmode)
}
