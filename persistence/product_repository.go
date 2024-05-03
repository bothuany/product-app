package persistence

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/gommon/log"
	"product-app/domain"
	"product-app/persistence/common"
)

type IProductRepository interface {
	GetAllProducts() []domain.Product
	GetAllProductsByStore(storeName string) []domain.Product
	AddProduct(product domain.Product) error
	GetProductById(productId int64) (domain.Product, error)
	DeleteProductById(productId int64) error
	UpdateProduct(product domain.Product) error
}

type ProductRepository struct {
	dbPool *pgxpool.Pool
}

func NewProductRepository(dbPool *pgxpool.Pool) IProductRepository {
	return &ProductRepository{dbPool: dbPool}
}

func (productRepository *ProductRepository) GetAllProducts() []domain.Product {
	ctx := context.Background()
	productRows, err := productRepository.dbPool.Query(ctx, "SELECT * FROM products")

	if err != nil {
		log.Error("Error while fetching products: %v\n", err)
		return []domain.Product{}
	}

	return extractProductsForRows(productRows)
}

func (productRepository *ProductRepository) GetAllProductsByStore(storeName string) []domain.Product {
	ctx := context.Background()

	getProductsByStoreNameSql := `SELECT * FROM products WHERE store = $1`

	productRows, err := productRepository.dbPool.Query(ctx, getProductsByStoreNameSql, storeName)

	if err != nil {
		log.Error("Error while fetching products: %v\n", err)
		return []domain.Product{}
	}

	return extractProductsForRows(productRows)
}

func (productRepository *ProductRepository) AddProduct(product domain.Product) error {
	ctx := context.Background()

	addProductSql := `INSERT INTO products (name, price, discount, store) VALUES ($1, $2, $3, $4)`

	addNewProduct, err := productRepository.dbPool.Exec(ctx, addProductSql, product.Name, product.Price, product.Discount, product.Store)

	if err != nil {
		log.Error("Error while adding product: %v\n", err)
		return err
	}
	log.Info("Product added with id: %v\n", addNewProduct)

	return nil
}

func (productRepository *ProductRepository) GetProductById(productId int64) (domain.Product, error) {
	ctx := context.Background()

	getProductByIdSql := `SELECT * FROM products WHERE id = $1`

	queryRow := productRepository.dbPool.QueryRow(ctx, getProductByIdSql, productId)

	var id int64
	var name string
	var price float32
	var discount float32
	var store string

	err := queryRow.Scan(&id, &name, &price, &discount, &store)

	if err != nil && err.Error() == common.NOT_FOUND {
		return domain.Product{}, errors.New(fmt.Sprintf("Product with id %d not found", productId))
	}
	if err != nil {
		log.Error("Error while fetching product: %v\n", err)
		return domain.Product{}, errors.New(fmt.Sprintf("Error while getting product with id %d not found", productId))
	}

	log.Info("Product fetched with id: %v\n", id)
	return domain.Product{Id: id, Name: name, Price: price, Discount: discount, Store: store}, nil

}

func (productRepository *ProductRepository) DeleteProductById(productId int64) error {
	ctx := context.Background()

	getErr := productRepository.checkProductExists(productId)
	if getErr != nil {
		return getErr
	}

	deleteProductSql := `DELETE FROM products WHERE id = $1`

	_, err := productRepository.dbPool.Exec(ctx, deleteProductSql, productId)

	if err != nil {
		log.Error("Error while deleting product: %v\n", err)
		return errors.New(fmt.Sprintf("Error while deleting product with id %d", productId))
	}

	log.Info("Product deleted with id: %v\n", productId)
	return nil
}

func (productRepository *ProductRepository) UpdateProduct(product domain.Product) error {
	ctx := context.Background()

	getErr := productRepository.checkProductExists(product.Id)
	if getErr != nil {
		return getErr
	}

	updateProductSql := `UPDATE products SET `
	var updateArgs []interface{} // Slice to hold the arguments for the query, interface{} is used to hold any type
	var argId = 1

	if product.Name != "" {
		updateProductSql += fmt.Sprintf("name = $%d,", argId)
		updateArgs = append(updateArgs, product.Name)
		argId++
	}
	if product.Price != 0 {
		updateProductSql += fmt.Sprintf("price = $%d,", argId)
		updateArgs = append(updateArgs, product.Price)
		argId++
	}
	if product.Discount != 0 {
		updateProductSql += fmt.Sprintf("discount = $%d,", argId)
		updateArgs = append(updateArgs, product.Discount)
		argId++
	}
	if product.Store != "" {
		updateProductSql += fmt.Sprintf("store = $%d,", argId)
		updateArgs = append(updateArgs, product.Store)
		argId++
	}

	// Remove the trailing comma and add the WHERE clause
	updateProductSql = updateProductSql[:len(updateProductSql)-1] + fmt.Sprintf(" WHERE id = $%d", argId)
	updateArgs = append(updateArgs, product.Id)

	_, err := productRepository.dbPool.Exec(ctx, updateProductSql, updateArgs...)

	if err != nil {
		log.Error("Error while updating product: %v\n", err)
		return errors.New(fmt.Sprintf("Error while updating product with id %d", product.Id))
	}

	log.Info("Product updated with id: %v\n", product.Id)
	return nil
}

func extractProductsForRows(productRows pgx.Rows) []domain.Product {
	var products []domain.Product

	var id int64
	var name string
	var price float32
	var discount float32
	var store string

	for productRows.Next() {
		err := productRows.Scan(&id, &name, &price, &discount, &store)
		if err != nil {
			log.Error("Error while scanning product: %v\n", err)
			return []domain.Product{}
		}
		products = append(products, domain.Product{Id: id, Name: name, Price: price, Discount: discount, Store: store})
	}

	log.Info("Products fetched: %v\n", products)
	return products
}

func (productRepository *ProductRepository) checkProductExists(productId int64) error {
	_, getErr := productRepository.GetProductById(productId)

	if getErr != nil {
		log.Error("Error while checking product: %v\n", getErr)
		return errors.New(fmt.Sprintf("Product with id %d not found", productId))
	}

	return nil
}
