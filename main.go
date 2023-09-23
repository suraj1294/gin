package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/suraj1294/go-gin-planetscale/handler"
)

var db *sql.DB

type Product struct {
	Id    int64
	Name  string
	Price int
}

func main() {
	// Load in the `.env` file

	productHandler := handler.GetProductHandler()

	// Build router & define routes
	router := gin.Default()
	router.GET("/products", productHandler.GetProducts)
	router.GET("/products/:productId", productHandler.GetProduct)
	router.POST("/products", productHandler.AddProduct)
	router.PUT("/products/:productId", productHandler.UpdateProduct)
	router.DELETE("/products/:productId", DeleteProduct)

	// Run the router
	router.Run()

}

func UpdateProduct(c *gin.Context) {
	var updates Product
	err := c.BindJSON(&updates)
	if err != nil {
		log.Fatal("(UpdateProduct) c.BindJSON", err)
	}

	productId := c.Param("productId")
	productId = strings.ReplaceAll(productId, "/", "")
	productIdInt, err := strconv.Atoi(productId)
	if err != nil {
		log.Fatal("(UpdateProduct) strconv.Atoi", err)
	}

	query := `UPDATE products SET name = ?, price = ? WHERE id = ?`
	_, err = db.Exec(query, updates.Name, updates.Price, productIdInt)
	if err != nil {
		log.Fatal("(UpdateProduct) db.Exec", err)
	}

	c.Status(http.StatusOK)
}

func DeleteProduct(c *gin.Context) {
	productId := c.Param("productId")

	productId = strings.ReplaceAll(productId, "/", "")
	productIdInt, err := strconv.Atoi(productId)
	if err != nil {
		log.Fatal("(DeleteProduct) strconv.Atoi", err)
	}
	query := `DELETE FROM products WHERE id = ?`
	_, err = db.Exec(query, productIdInt)
	if err != nil {
		log.Fatal("(DeleteProduct) db.Exec", err)
	}

	c.Status(http.StatusOK)
}
