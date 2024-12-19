package products

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Delete(db *sql.DB) gin.HandlerFunc {
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

		var isProductInTransactions bool
		sqlCheckProductInTransactions := "SELECT EXISTS (SELECT 1 FROM bill_details WHERE product_id = $1)"
		if err := db.QueryRow(sqlCheckProductInTransactions, id).Scan(&isProductInTransactions); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check product in transaction existence"})
			return
		}
		if isProductInTransactions {
			c.JSON(http.StatusConflict, gin.H{"error": "Product id is being used in transcations"})
			return
		}

		sqlDeleteProduct := "DELETE FROM products WHERE id = $1"
		_, err := db.Exec(sqlDeleteProduct, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully", "data": "OK"})
	}
}