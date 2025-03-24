package prodcuts

type ProductService struct {
	repo ProductRepo
}

func NewProductService(repo ProductRepo) *ProductService {
	return &ProductService{repo}
}

func (s *ProductService) CreateProduct(name, description string, price float64, category string, stock int) (int, error) {
	return s.repo.CreateProduct(name, description, price, category, stock)
}
