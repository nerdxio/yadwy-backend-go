package prodcuts

type Product struct {
	id          int
	name        string
	description string
	price       float64
	categoryId  string
	stock       int
	createdAt   string
	updatedAt   string
}

type ProductParams struct {
	id          int
	name        string
	description string
	price       float64
	categoryId  string
	stock       int
	createdAt   string
	updatedAt   string
}

func NewProduct(params ProductParams) *Product {
	return &Product{
		id:          params.id,
		name:        params.name,
		description: params.description,
		price:       params.price,
		categoryId:  params.categoryId,
		stock:       params.stock,
		createdAt:   params.createdAt,
		updatedAt:   params.updatedAt,
	}
}
