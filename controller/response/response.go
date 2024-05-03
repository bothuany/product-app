package response

import "product-app/domain"

type ErrorResponse struct {
	ErrorDescription string `json:"errorDescription"`
}

type GetProductByIdResponse struct {
	Id       int64   `json:"id"`
	Name     string  `json:"name"`
	Price    float32 `json:"price"`
	Discount float32 `json:"discount"`
	Store    string  `json:"store"`
}

func (productResponse GetProductByIdResponse) ToResponse(product *domain.Product) GetProductByIdResponse {
	return GetProductByIdResponse{
		Id:       product.Id,
		Name:     product.Name,
		Price:    product.Price,
		Discount: product.Discount,
		Store:    product.Store,
	}
}

type GetAllProductsResponse struct {
	Id       int64   `json:"id"`
	Name     string  `json:"name"`
	Price    float32 `json:"price"`
	Discount float32 `json:"discount"`
	Store    string  `json:"store"`
}

func (productResponse GetAllProductsResponse) ToResponse(product *domain.Product) GetAllProductsResponse {
	return GetAllProductsResponse{
		Id:       product.Id,
		Name:     product.Name,
		Price:    product.Price,
		Discount: product.Discount,
		Store:    product.Store,
	}
}
