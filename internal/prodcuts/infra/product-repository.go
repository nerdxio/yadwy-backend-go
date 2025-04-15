package infra

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strings"
	"time"
	"yadwy-backend/internal/prodcuts/domain"
)

type ProductRepositoryImpl struct {
	db *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) domain.ProductRepository {
	return &ProductRepositoryImpl{db: db}
}

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

	if len(p.Labels) > 0 {
		for _, label := range p.Labels {
			query := `INSERT INTO product_labels (product_id, label_name)
                     VALUES ($1, $2)`
			_, err = tx.ExecContext(ctx, query, p.ID, label)
			if err != nil {
				tx.Rollback()
				return err
			}
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

func (r *ProductRepositoryImpl) SearchProducts(ctx context.Context, params domain.SearchParams) (*domain.SearchResult, error) {
	// Start building the query
	query := "SELECT DISTINCT p.* FROM products p"
	countQuery := "SELECT COUNT(DISTINCT p.id) FROM products p"

	// Initialize where clauses and arguments
	whereClauses := []string{}
	args := []interface{}{}
	argIndex := 1

	// Add join for labels if needed
	if len(params.Labels) > 0 {
		query += " LEFT JOIN product_labels pl ON p.id = pl.product_id"
		countQuery += " LEFT JOIN product_labels pl ON p.id = pl.product_id"
	}

	// Add search conditions
	if params.Query != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("(p.name ILIKE $%d OR p.description ILIKE $%d)", argIndex, argIndex))
		args = append(args, "%"+params.Query+"%")
		argIndex++
	}

	if params.CategoryID != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("p.category_id = $%d", argIndex))
		args = append(args, params.CategoryID)
		argIndex++
	}

	if params.MinPrice != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("p.price >= $%d", argIndex))
		args = append(args, *params.MinPrice)
		argIndex++
	}

	if params.MaxPrice != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("p.price <= $%d", argIndex))
		args = append(args, *params.MaxPrice)
		argIndex++
	}

	if params.SellerID != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("p.seller_id = $%d", argIndex))
		args = append(args, *params.SellerID)
		argIndex++
	}

	if params.Available != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("p.is_available = $%d", argIndex))
		args = append(args, *params.Available)
		argIndex++
	}

	// Handle labels filtering
	if len(params.Labels) > 0 {
		placeholders := make([]string, len(params.Labels))
		for i := range params.Labels {
			placeholders[i] = fmt.Sprintf("$%d", argIndex)
			args = append(args, params.Labels[i])
			argIndex++
		}
		whereClauses = append(whereClauses, fmt.Sprintf("pl.label_name IN (%s)", strings.Join(placeholders, ", ")))
	}

	// Add WHERE clause if we have conditions
	if len(whereClauses) > 0 {
		whereClause := " WHERE " + strings.Join(whereClauses, " AND ")
		query += whereClause
		countQuery += whereClause
	}

	// Add sorting
	if params.SortBy != "" {
		allowedSortFields := map[string]bool{
			"price": true, "created_at": true, "name": true,
		}
		if allowedSortFields[params.SortBy] {
			sortDir := "ASC"
			if params.SortDir == "desc" {
				sortDir = "DESC"
			}
			query += fmt.Sprintf(" ORDER BY p.%s %s", params.SortBy, sortDir)
		}
	} else {
		query += " ORDER BY p.created_at DESC"
	}

	// Add pagination
	limit := 10 // Default limit
	if params.Limit > 0 {
		limit = params.Limit
	}

	offset := 0 // Default offset
	if params.Offset > 0 {
		offset = params.Offset
	}

	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, limit, offset)

	// Execute count query first
	var totalCount int
	err := r.db.GetContext(ctx, &totalCount, countQuery, args[:argIndex-1]...)
	if err != nil {
		return nil, fmt.Errorf("failed to count search results: %w", err)
	}

	// Execute main query
	rows, err := r.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search products: %w", err)
	}
	defer rows.Close()

	// Parse results
	products := []*domain.Product{}
	for rows.Next() {
		var pdb productDB
		if err := rows.StructScan(&pdb); err != nil {
			return nil, fmt.Errorf("failed to scan product row: %w", err)
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

		products = append(products, product)
	}

	// Load images for each product
	for _, product := range products {
		imgRows, err := r.db.QueryxContext(ctx,
			"SELECT * FROM product_images WHERE product_id = $1", product.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch product images: %w", err)
		}

		for imgRows.Next() {
			var img imageDB
			if err := imgRows.StructScan(&img); err != nil {
				imgRows.Close()
				return nil, fmt.Errorf("failed to scan image row: %w", err)
			}
			product.Images = append(product.Images, domain.Image{
				URL:  img.ImageURL,
				Type: img.ImageType,
			})
		}
		imgRows.Close()

		// Fetch labels for this product
		labelRows, err := r.db.QueryxContext(ctx,
			"SELECT label_name FROM product_labels WHERE product_id = $1", product.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch product labels: %w", err)
		}

		for labelRows.Next() {
			var label string
			if err := labelRows.Scan(&label); err != nil {
				labelRows.Close()
				return nil, fmt.Errorf("failed to scan label row: %w", err)
			}
			product.Labels = append(product.Labels, label)
		}
		labelRows.Close()
	}

	result := &domain.SearchResult{
		Products:    products,
		TotalCount:  totalCount,
		Limit:       limit,
		Offset:      offset,
		HasNextPage: (offset + len(products)) < totalCount,
	}

	return result, nil
}
