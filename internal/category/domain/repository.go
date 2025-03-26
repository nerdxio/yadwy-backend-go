package domain

// CategoryRepository defines the interface for category data operations
type CategoryRepository interface {
	CreateCategory(name, description string) (int, error)
	GetCategory(id int) (*Category, error)
	ListCategories() ([]Category, error)
	UpdateCategory(id int, name, description string) error
	DeleteCategory(id int) error
}
