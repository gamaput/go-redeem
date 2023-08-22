package repository

import (
	"github.com/gamaput/go-redeem/model"
	"gorm.io/gorm"
)

type productRepository struct {
	DB *gorm.DB
}

// ProductRepository : represent the product's repository contract
type ProductRepository interface {
	CreateProduct(model.Product) (model.Product, error)
	GetProductByID(int) (model.Product, error)
	UpdateProduct(model.Product) (model.Product, error)
	DeleteProduct(model.Product) (model.Product, error)
	GetAllProducts() ([]model.Product, error)
}

// NewProductRepository -> returns new product repository
func NewProductRepository(db *gorm.DB) ProductRepository {
	return productRepository{
		DB: db,
	}
}

func (pr productRepository) CreateProduct(product model.Product) (model.Product, error) {
	return product, pr.DB.Create(&product).Error
}

func (pr productRepository) GetProductByID(id int) (product model.Product, err error) {
	return product, pr.DB.First(&product, id).Error
}

func (pr productRepository) UpdateProduct(product model.Product) (model.Product, error) {
	if err := pr.DB.Model(&model.Product{}).Where("id =?", product.ID).Updates(&product).Error; err != nil {
		return product, err
	}
	return product, nil
}

func (pr productRepository) DeleteProduct(product model.Product) (model.Product, error) {
	if err := pr.DB.First(&product, product.ID).Error; err != nil {
		return product, err
	}
	return product, pr.DB.Delete(&product).Error
}

func (pr productRepository) GetAllProducts() (products []model.Product, err error) {
	return products, pr.DB.Find(&products).Error
}
