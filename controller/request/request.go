package request

import "product-app/service/model"

type AddProductRequest struct {
	Name     string  `json:"name"`
	Price    float32 `json:"price"`
	Discount float32 `json:"discount"`
	Store    string  `json:"store"`
}

func (addProductRequest AddProductRequest) ToModel() model.ProductCreate {
	return model.ProductCreate{
		Name:     addProductRequest.Name,
		Price:    addProductRequest.Price,
		Discount: addProductRequest.Discount,
		Store:    addProductRequest.Store,
	}
}

type UpdateProductRequest struct {
	Id       int64   `json:"id"`
	Name     string  `json:"name"`
	Price    float32 `json:"price"`
	Discount float32 `json:"discount"`
	Store    string  `json:"store"`
}

func (updateProductRequest UpdateProductRequest) ToModel() model.ProductUpdate {
	return model.ProductUpdate{
		Id:       updateProductRequest.Id,
		Name:     updateProductRequest.Name,
		Price:    updateProductRequest.Price,
		Discount: updateProductRequest.Discount,
		Store:    updateProductRequest.Store,
	}
}
