package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/suraj1294/go-gin-planetscale/logger"
	"github.com/suraj1294/go-gin-planetscale/services"
)

type ProductHandler struct {
	services *services.ProductService
}

func (product *ProductHandler) GetProducts(c *gin.Context) {

	products, err := product.services.GetAll()

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}

	c.JSON(http.StatusOK, products)

}

func (product *ProductHandler) GetProduct(c *gin.Context) {

	productId := c.Param("productId")
	productId = strings.ReplaceAll(productId, "/", "")
	productIdInt, err := strconv.Atoi(productId)
	if err != nil {
		logger.Error("(GetSingleProduct) strconv.Atoi" + err.Error())
	}

	products, selectErr := product.services.GetById(productIdInt)

	if selectErr != nil {
		c.JSON(http.StatusNotFound, "product not found")
	}

	c.JSON(http.StatusOK, products)

}

func (product *ProductHandler) AddProduct(c *gin.Context) {

	var newProduct *services.Product
	err := c.BindJSON(&newProduct)

	if err != nil {
		c.JSON(http.StatusInternalServerError, "failed to add product")
	}

	products, selectErr := product.services.Add(newProduct)

	if selectErr != nil {
		c.JSON(http.StatusNotFound, "product not found")
	}

	c.JSON(http.StatusOK, products)

}

func (product *ProductHandler) UpdateProduct(c *gin.Context) {

	var updates *services.Product
	err := c.Bind(&updates)

	if err != nil {
		logger.Error("(UpdateProduct) c.BindJSON " + err.Error())

		c.JSON(http.StatusBadRequest, gin.H{"errros": fmt.Sprintf("%v", err)})
	}

	productId := c.Param("productId")
	productId = strings.ReplaceAll(productId, "/", "")
	productIdInt, err := strconv.Atoi(productId)

	if err != nil {
		logger.Error("(UpdateProduct) strconv.Atoi " + err.Error())
	}

	products, updateErr := product.services.Update(updates, productIdInt)

	if updateErr != nil {
		c.JSON(http.StatusInternalServerError, "failed to update product")
	}

	c.JSON(http.StatusOK, products)

}

func GetProductHandler() *ProductHandler {
	return &ProductHandler{services: services.GetProductService()}
}
