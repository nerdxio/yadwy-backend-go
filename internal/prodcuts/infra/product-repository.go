package infra

import (
	"context"
	"github.com/jmoiron/sqlx"
	"time"
	"yadwy-backend/internal/prodcuts/domain"
)

type ProductRepositoryImpl struct {
	db *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) domain.ProductRepository {
	return &ProductRepositoryImpl{db: db}
}

// Database models
type productDB struct {
	ID          int64   `db:"id"`
	Name        string  `db:"name"`
	Description string  `db:"description"`
	Price       float64 `db:"price"`
	CategoryID  string  `db:"category_id"`
	SellerID    int64   `db:"seller_id"`
	Stock       int     `db:"stock"`
	IsAvailable bool    `db:"is_available"`
	CreatedAt   string  `db:"created_at"`
	UpdatedAt   string  `db:"updated_at"`
}

type imageDB struct {
	ID        int64     `db:"id"`
	ProductID int64     `db:"product_id"`
	ImageURL  string    `db:"image_url"`
	ImageType string    `db:"image_type"`
	CreatedAt time.Time `db:"created_at"`
}

func (r *ProductRepositoryImpl) CreateProduct(ctx context.Context, p *domain.Product, images []domain.Image) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	query := `INSERT INTO products (name, description, price, category_id, seller_id, stock, is_available)
              VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`

	err = tx.QueryRowxContext(ctx, query, p.Name, p.Description, p.Price, p.CategoryID,
		p.SellerID, p.Stock, p.IsAvailable).Scan(&p.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, img := range images {
		query := `INSERT INTO product_images (product_id, image_url, image_type)
                 VALUES ($1, $2, $3)`
		_, err = tx.ExecContext(ctx, query, p.ID, img.URL, img.Type)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (r *ProductRepositoryImpl) GetProduct(ctx context.Context, id int64) (*domain.Product, error) {
	var pdb productDB
	err := r.db.GetContext(ctx, &pdb,
		"SELECT * FROM products WHERE id = $1", id)
	if err != nil {
		return nil, err
	}

	product := &domain.Product{
		ID:          pdb.ID,
		Name:        pdb.Name,
		Description: pdb.Description,
		Price:       pdb.Price,
		CategoryID:  pdb.CategoryID,
		SellerID:    pdb.SellerID,
		Stock:       pdb.Stock,
		IsAvailable: pdb.IsAvailable,
		CreatedAt:   pdb.CreatedAt,
		UpdatedAt:   pdb.UpdatedAt,
	}

	rows, err := r.db.QueryxContext(ctx,
		"SELECT * FROM product_images WHERE product_id = $1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var img imageDB
		if err := rows.StructScan(&img); err != nil {
			return nil, err
		}
		product.Images = append(product.Images, domain.Image{
			URL:  img.ImageURL,
			Type: img.ImageType,
		})
	}

	return product, nil
}
