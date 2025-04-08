package category

type CategoryRepo interface {
	CreateCategory(name, description, imageUrl string) (int, error)
}
