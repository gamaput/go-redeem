// controller/product_controller.go

package controller

import (
	"net/http"
	"strconv"

	"github.com/casbin/casbin/v2"
	"github.com/gamaput/go-redeem/model"
	"github.com/gamaput/go-redeem/repository"
	"github.com/gin-gonic/gin"
)

// ProductController : represent the product's controller contract
type ProductController interface {
	CreateProduct(enforcer *casbin.Enforcer) gin.HandlerFunc
	GetProductByID(*gin.Context)
	UpdateProduct(*gin.Context)
	DeleteProduct(*gin.Context)
	GetAllProducts(*gin.Context)
}

type productController struct {
	productRepo repository.ProductRepository
}

// NewProductController -> returns new product controller
func NewProductController(productRepo repository.ProductRepository) ProductController {
	return productController{
		productRepo: productRepo,
	}
}

func (pc productController) CreateProduct(enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var product model.Product
		if err := ctx.ShouldBindJSON(&product); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		product, err := pc.productRepo.CreateProduct(product)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, product)
	}
}

func (pc productController) GetProductByID(c *gin.Context) {
	id := c.Param("product")
	intID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	product, err := pc.productRepo.GetProductByID(intID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (pc productController) UpdateProduct(c *gin.Context) {
	var product model.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id := c.Param("product")
	intID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	product.ID = uint(intID)
	product, err = pc.productRepo.UpdateProduct(product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, product)
}

func (pc productController) DeleteProduct(c *gin.Context) {
	var product model.Product
	id := c.Param("product")
	intID, _ := strconv.Atoi(id)
	product.ID = uint(intID)
	product, err := pc.productRepo.DeleteProduct(product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, product)
}

func (pc productController) GetAllProducts(c *gin.Context) {

	products, err := pc.productRepo.GetAllProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}
