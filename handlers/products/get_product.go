package products

import (
	"challenge-goapi/models"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Get(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var product models.Product

		if id != "" {
			sqlGetProductById := "SELECT id, name, price, unit FROM products WHERE id = $1"
			if err := db.QueryRow(sqlGetProductById, id).Scan(&product.Id, &product.Name, &product.Price, &product.Unit); err != nil {
				if err == sql.ErrNoRows {
					c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
				} else {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
				}
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": "Product retrieved successfully", "data": product})
			return
		}

		sqlGetAllProduct := "SELECT id, name, price, unit FROM products"
		rows, err := db.Query(sqlGetAllProduct)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Failed to retrieve products from database"})
			return
		}
		defer rows.Close()

		var listProducts []models.Product
		for rows.Next() {
			if err := rows.Scan(&product.Id, &product.Name, &product.Price, &product.Unit); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan product data"})
				return
			}
			listProducts = append(listProducts, product)
		}
		if err := rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while retrieving products"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "List of products retrieved successfully","data": listProducts})
	}
}
