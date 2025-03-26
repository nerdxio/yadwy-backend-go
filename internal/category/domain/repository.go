package domain

// CategoryRepository defines the data access interface for categories
type CategoryRepository interface {
	CreateCategory(name, description string) (int, error)
	// Add other methods as needed
}
