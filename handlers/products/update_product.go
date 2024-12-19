package products

import (
	"challenge-goapi/models"
	"database/sql"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Put(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var idIsExsits bool
		sqlCheckProduct := "SELECT EXISTS (SELECT 1 FROM products WHERE id = $1);"
		if err := db.QueryRow(sqlCheckProduct, id).Scan(&idIsExsits); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check product existence"})
			return
		}
		if !idIsExsits {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}

		var requestProduct models.ProductUpdate
		if err := c.ShouldBindJSON(&requestProduct); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
			return
		}

		var currentProduct models.Product
		sqlSelectProduct := "SELECT id, name, price, unit FROM products WHERE id = $1"
		if err := db.QueryRow(sqlSelectProduct, id).Scan(&currentProduct.Id, &currentProduct.Name, &currentProduct.Price, &currentProduct.Unit); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve current product data"})
			return
		}

		if requestProduct.Name != nil && strings.TrimSpace(*requestProduct.Name) != "" {
			currentProduct.Name = *requestProduct.Name
		}
		if requestProduct.Price != nil && *requestProduct.Price > 0 {
			currentProduct.Price = *requestProduct.Price
		}
		if requestProduct.Unit != nil && strings.TrimSpace(*requestProduct.Unit) != "" {
			currentProduct.Name = *requestProduct.Name
		}

		sqlUpdateProduct := "UPDATE products SET name = $2, price = $3, unit = $4 WHERE id = $1"
		_, err := db.Exec(sqlUpdateProduct, currentProduct.Id, currentProduct.Name, currentProduct.Price, currentProduct.Unit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update prdouct"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Customer updated successfully", "data": currentProduct})
	}
}