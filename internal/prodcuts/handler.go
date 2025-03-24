package prodcuts

type ProductController struct {
	service ProductService
}

func NewProductHandler(service ProductService) *ProductController {
	return &ProductController{service}
}
