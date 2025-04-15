package domain

import "context"

// SearchParams contains all possible search parameters
type SearchParams struct {
	Query      string   // Text search in name/description
	CategoryID string   // Filter by category
	MinPrice   *float64 // Minimum price
	MaxPrice   *float64 // Maximum price
	Labels     []string // Filter by labels
	SellerID   *int64   // Filter by seller
	Available  *bool    // Filter by availability
	Limit      int      // Pagination limit
	Offset     int      // Pagination offset
	SortBy     string   // Field to sort by
	SortDir    string   // Sort direction (asc/desc)
}

// SearchResult represents paginated search results
type SearchResult struct {
	Products    []*Product // List of products
	TotalCount  int        // Total count of matching products
	Limit       int        // Items per page
	Offset      int        // Current page offset
	HasNextPage bool       // Whether there are more results
}

// ProductRepository interface extension
type ProductRepository interface {
	CreateProduct(ctx context.Context, p *Product, images []Image) error
	GetProduct(ctx context.Context, id int64) (*Product, error)
	SearchProducts(ctx context.Context, params SearchParams) (*SearchResult, error)
}
