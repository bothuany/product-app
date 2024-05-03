package infrastructure

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
	"os"
	"product-app/common/postgresql"
	"product-app/domain"
	"product-app/persistence"
	"testing"
)

var productRepository persistence.IProductRepository
var dbPool *pgxpool.Pool
var ctx context.Context

func TestMain(m *testing.M) {
	ctx = context.Background()

	dbPool = postgresql.GetConnectionPool(ctx, postgresql.Config{
		Host:                  "localhost",
		Port:                  "6432",
		DbName:                "productapp",
		UserName:              "postgres",
		Password:              "postgres",
		MaxConnections:        "10",
		MaxConnectionIdleTime: "30s",
	})

	productRepository = persistence.NewProductRepository(dbPool)
	fmt.Println("-------------Before all tests-------------")
	exitCode := m.Run()
	fmt.Println("-------------After all tests-------------")
	os.Exit(exitCode)
}

func setup(ctx context.Context, dbPool *pgxpool.Pool) {
	TestDataInitialize(ctx, dbPool)
}

func clear(ctx context.Context, dbPool *pgxpool.Pool) {
	TruncateTestData(ctx, dbPool)
}
func TestGetAllProducts(t *testing.T) {
	setup(ctx, dbPool)
	expectedProducts := []domain.Product{
		{Id: 1, Name: "AirFryer", Price: 3000.0, Discount: 22.0, Store: "ABC TECH"},
		{Id: 2, Name: "Ütü", Price: 1500.0, Discount: 10.0, Store: "ABC TECH"},
		{Id: 3, Name: "Çamaşır Makinesi", Price: 10000.0, Discount: 15.0, Store: "ABC TECH"},
		{Id: 4, Name: "Lambader", Price: 2000.0, Discount: 0.0, Store: "Dekorasyon Sarayı"},
	}

	t.Run("GetAllProducts", func(t *testing.T) {
		actualProducts := productRepository.GetAllProducts()
		assert.Equal(t, 4, len(actualProducts))
		assert.Equal(t, expectedProducts, actualProducts)
	})

	clear(ctx, dbPool)

}

func TestGetAllProductsByStore(t *testing.T) {
	setup(ctx, dbPool)
	expectedProducts := []domain.Product{
		{Id: 1, Name: "AirFryer", Price: 3000.0, Discount: 22.0, Store: "ABC TECH"},
		{Id: 2, Name: "Ütü", Price: 1500.0, Discount: 10.0, Store: "ABC TECH"},
		{Id: 3, Name: "Çamaşır Makinesi", Price: 10000.0, Discount: 15.0, Store: "ABC TECH"},
	}

	t.Run("GetAllProductsByStore", func(t *testing.T) {
		actualProducts := productRepository.GetAllProductsByStore("ABC TECH")
		assert.Equal(t, 3, len(actualProducts))
		assert.Equal(t, expectedProducts, actualProducts)
	})

	clear(ctx, dbPool)
}

func TestAddProduct(t *testing.T) {

	expectedProducts := []domain.Product{
		{Id: 1, Name: "Kupa", Price: 100.0, Discount: 0.0, Store: "RBD"},
	}

	newProduct := domain.Product{
		Name:     "Kupa",
		Price:    100.0,
		Discount: 0.0,
		Store:    "RBD",
	}
	t.Run("AddProduct", func(t *testing.T) {
		productRepository.AddProduct(newProduct)
		actualProducts := productRepository.GetAllProducts()
		assert.Equal(t, 1, len(actualProducts))
		assert.Equal(t, expectedProducts, actualProducts)
	})

	clear(ctx, dbPool)
}

func TestGetProductById(t *testing.T) {
	setup(ctx, dbPool)
	expectedProduct := domain.Product{Id: 1, Name: "AirFryer", Price: 3000.0, Discount: 22.0, Store: "ABC TECH"}
	t.Run("GetProductById", func(t *testing.T) {
		actualProduct, _ := productRepository.GetProductById(1)
		_, err := productRepository.GetProductById(5)
		assert.Equal(t, expectedProduct, actualProduct)
		assert.Equal(t, "Product with id 5 not found", err.Error())
	})
	clear(ctx, dbPool)
}

func TestDeleteProductById(t *testing.T) {
	setup(ctx, dbPool)

	t.Run("DeleteProductById", func(t *testing.T) {
		productRepository.DeleteProductById(4)
		_, err := productRepository.GetProductById(4)
		assert.Equal(t, "Product with id 4 not found", err.Error())
	})

	clear(ctx, dbPool)
}

func TestUpdateProduct(t *testing.T) {
	setup(ctx, dbPool)
	expectedProduct := domain.Product{Id: 1, Name: "Fırın", Price: 4000.0, Discount: 22.0, Store: "ABC TECH"}
	updatedProduct := domain.Product{Id: 1, Name: "Fırın", Price: 4000.0, Discount: 22.0, Store: "ABC TECH"}

	t.Run("UpdateProduct", func(t *testing.T) {
		productBeforeUpdate, _ := productRepository.GetProductById(1)
		assert.Equal(t, "AirFryer", productBeforeUpdate.Name)
		assert.Equal(t, float32(3000.0), productBeforeUpdate.Price)
		productRepository.UpdateProduct(updatedProduct)
		productAfterUpdate, _ := productRepository.GetProductById(1)
		assert.Equal(t, expectedProduct, productAfterUpdate)
	})

	clear(ctx, dbPool)
}
