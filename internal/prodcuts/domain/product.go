package domain

type Product struct {
	ID          int64    `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       float64  `json:"price"`
	CategoryID  string   `json:"category_id"`
	SellerID    int64    `json:"seller_id"`
	Stock       int      `json:"stock"`
	IsAvailable bool     `json:"is_available"`
	Images      []Image  `json:"images"`
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
	Labels      []string `json:"labels"`
}

type Image struct {
	URL  string `json:"url"`
	Type string `json:"type"` // "thumbnail", "main", "extra"
}
