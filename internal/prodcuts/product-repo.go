package prodcuts

type ProductRepo interface {
	CreateProduct(name, description string, price float64, category string, stock int) (int, error)
}
