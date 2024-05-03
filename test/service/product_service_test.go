package service

import (
	"github.com/stretchr/testify/assert"
	"os"
	"product-app/domain"
	"product-app/service"
	"product-app/service/model"
	"testing"
)

var productService service.IProductService

func TestMain(m *testing.M) {
	initialProducts := []domain.Product{
		{Id: 1, Name: "AirFryer", Price: 3000.0, Discount: 22.0, Store: "ABC TECH"},
		{Id: 2, Name: "Ütü", Price: 1500.0, Discount: 10.0, Store: "ABC TECH"},
		{Id: 3, Name: "Çamaşır Makinesi", Price: 10000.0, Discount: 15.0, Store: "ABC TECH"},
		{Id: 4, Name: "Lambader", Price: 2000.0, Discount: 0.0, Store: "Dekorasyon Sarayı"},
	}
	fakeProductRepository := NewFakeProductRepository(initialProducts)
	productService = service.NewProductService(fakeProductRepository)
	exitCode := m.Run()
	os.Exit(exitCode)
}

func Test_ShouldGetAllProducts(t *testing.T) {
	t.Run("ShouldGetAllProducts", func(t *testing.T) {
		products := productService.GetAllProducts()
		assert.Equal(t, 4, len(products))
	})
}

func Test_ShouldGetAllProductsByStore(t *testing.T) {
	t.Run("ShouldGetAllProductsByStore", func(t *testing.T) {
		products := productService.GetAllProductsByStore("ABC TECH")
		assert.Equal(t, 3, len(products))
	})
}

func Test_WhenNoValidationErrorOccurred_ShouldAddProduct(t *testing.T) {
	t.Run("WhenNoValidationErrorOccurred_ShouldAddProduct", func(t *testing.T) {
		productService.AddProduct(model.ProductCreate{
			Name:     "Kupa",
			Price:    100.0,
			Discount: 0.0,
			Store:    "RBD",
		})

		actualProducts := productService.GetAllProducts()

		assert.Equal(t, 5, len(actualProducts))
		assert.Equal(t, "Kupa", actualProducts[len(actualProducts)-1].Name)
	})
}

func Test_WhenDiscountIsHigherThan70_ShouldNotAddProduct(t *testing.T) {
	t.Run("WhenDiscountIsHigherThan70_ShouldNotAddProduct", func(t *testing.T) {
		err := productService.AddProduct(model.ProductCreate{
			Name:     "Kupa",
			Price:    100.0,
			Discount: 71.0,
			Store:    "RBD",
		})

		assert.NotNil(t, err)
		assert.Equal(t, "Discount can not be greater than 70", err.Error())
		assert.Equal(t, 4, len(productService.GetAllProducts()))
	})
}

func Test_ShouldGetProductById(t *testing.T) {
	t.Run("ShouldGetProductById", func(t *testing.T) {
		product, err := productService.GetProductById(1)
		assert.Nil(t, err)
		assert.Equal(t, "AirFryer", product.Name)
	})
}

func Test_ShouldUpdateProduct(t *testing.T) {
	t.Run("ShouldUpdateProduct", func(t *testing.T) {
		err := productService.UpdateProduct(model.ProductUpdate{
			Id:       1,
			Name:     "Updated AirFryer",
			Price:    3500.0,
			Discount: 25.0,
			Store:    "ABC TECH",
		})
		assert.Nil(t, err)
		product, err := productService.GetProductById(1)
		assert.Nil(t, err)
		assert.Equal(t, "Updated AirFryer", product.Name)
		assert.Equal(t, float32(3500.0), product.Price)
		assert.Equal(t, float32(25.0), product.Discount)
	})
}

func Test_ShouldDeleteProductById(t *testing.T) {
	t.Run("ShouldDeleteProductById", func(t *testing.T) {
		err := productService.DeleteProductById(1)
		assert.Nil(t, err)
		_, err = productService.GetProductById(1)
		assert.NotNil(t, err)
	})
}
