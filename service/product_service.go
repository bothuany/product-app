package service

import (
	"errors"
	"product-app/domain"
	"product-app/persistence"
	"product-app/service/model"
)

type IProductService interface {
	AddProduct(productCreate model.ProductCreate) error
	DeleteProductById(productId int64) error
	GetProductById(productId int64) (domain.Product, error)
	UpdateProduct(productUpdate model.ProductUpdate) error
	GetAllProducts() []domain.Product
	GetAllProductsByStore(storeName string) []domain.Product
}

type ProductService struct {
	productRepository persistence.IProductRepository
}

func (productService ProductService) AddProduct(productCreate model.ProductCreate) error {
	validateErr := validateProductCreate(productCreate)

	if validateErr != nil {
		return validateErr
	}

	return productService.productRepository.AddProduct(domain.Product{
		Name:     productCreate.Name,
		Price:    productCreate.Price,
		Discount: productCreate.Discount,
		Store:    productCreate.Store,
	})
}

func (productService ProductService) DeleteProductById(productId int64) error {
	return productService.productRepository.DeleteProductById(productId)
}

func (productService ProductService) GetProductById(productId int64) (domain.Product, error) {
	product, err := productService.productRepository.GetProductById(productId)

	if err != nil {
		return domain.Product{}, err
	}

	return product, nil
}

func (productService ProductService) UpdateProduct(productUpdate model.ProductUpdate) error {
	return productService.productRepository.UpdateProduct(domain.Product{
		Id:       productUpdate.Id,
		Name:     productUpdate.Name,
		Price:    productUpdate.Price,
		Discount: productUpdate.Discount,
		Store:    productUpdate.Store,
	})
}

func (productService ProductService) GetAllProducts() []domain.Product {
	products := productService.productRepository.GetAllProducts()

	return products
}

func (productService ProductService) GetAllProductsByStore(storeName string) []domain.Product {
	products := productService.productRepository.GetAllProductsByStore(storeName)

	return products
}

func NewProductService(productRepository persistence.IProductRepository) IProductService {
	return &ProductService{productRepository: productRepository}
}

func validateProductCreate(productCreateRequest model.ProductCreate) error {
	if productCreateRequest.Discount > 70.0 {
		return errors.New("Discount can not be greater than 70")
	}

	return nil
}
