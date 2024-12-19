package customers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Delete(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		// check customer id ada atau ga
		var isCustomerIdExists bool
		sqlCheckCustomer := "SELECT EXISTS (SELECT 1 FROM customers WHERE id = $1);"
		if err := db.QueryRow(sqlCheckCustomer, id).Scan(&isCustomerIdExists); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check customer existence"})
			return
		}
		if !isCustomerIdExists {
			c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
			return
		}

		var isCustomerInTransactions bool
		sqlCheckCustomerInTransaction := "SELECT EXISTS (SELECT 1 FROM bills WHERE customer_id = $1)"
		if err := db.QueryRow(sqlCheckCustomerInTransaction, id).Scan(&isCustomerInTransactions); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check customer in transaction existence"})
			return
		}
		if isCustomerInTransactions {
			c.JSON(http.StatusConflict, gin.H{"error": "Customer id is being used in transcations"})
			return
		}

		sqlDeleteCustomer := "DELETE FROM customers WHERE id = $1;"
		_, err := db.Exec(sqlDeleteCustomer, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete customer"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Customer deleted successfully", "data": "OK"})
	}
}
