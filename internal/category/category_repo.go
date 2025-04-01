package category

type CategoryRepo interface {
	CreateCategory(name, description string) (int, error)
}
