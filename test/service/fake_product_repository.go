package service

import (
	"product-app/domain"
	"product-app/persistence"
)

type FakeProductRepository struct {
	products []domain.Product
}

func NewFakeProductRepository(initialProducts []domain.Product) persistence.IProductRepository {
	return &FakeProductRepository{
		products: initialProducts,
	}
}

func (fakeProductRepository *FakeProductRepository) GetAllProducts() []domain.Product {
	return fakeProductRepository.products
}

func (fakeProductRepository *FakeProductRepository) GetAllProductsByStore(storeName string) []domain.Product {
	var productsByStore []domain.Product

	for _, product := range fakeProductRepository.products {
		if product.Store == storeName {
			productsByStore = append(productsByStore, product)
		}
	}

	return productsByStore
}

func (fakeProductRepository *FakeProductRepository) AddProduct(product domain.Product) error {
	product.Id = int64(len(fakeProductRepository.products) + 1)
	fakeProductRepository.products = append(fakeProductRepository.products, product)
	return nil
}

func (fakeProductRepository *FakeProductRepository) GetProductById(productId int64) (domain.Product, error) {
	for _, product := range fakeProductRepository.products {
		if product.Id == productId {
			return product, nil
		}
	}

	return domain.Product{}, nil
}

func (fakeProductRepository *FakeProductRepository) DeleteProductById(productId int64) error {
	for i, product := range fakeProductRepository.products {
		if product.Id == productId {
			fakeProductRepository.products = append(fakeProductRepository.products[:i], fakeProductRepository.products[i+1:]...)
			return nil
		}
	}

	return nil
}

func (fakeProductRepository *FakeProductRepository) UpdateProduct(product domain.Product) error {
	for i, p := range fakeProductRepository.products {
		if p.Id == product.Id {
			fakeProductRepository.products[i] = product
			return nil
		}
	}

	return nil
}
