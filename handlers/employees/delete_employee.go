package employees

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Delete(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var isEmployeeIdExist bool
		sqlCheckEmployee := "SELECT EXISTS (SELECT 1 FROM employees WHERE id = $1);"
		if err := db.QueryRow(sqlCheckEmployee, id).Scan(&isEmployeeIdExist); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check employee existence"})
			return
		}
		if !isEmployeeIdExist {
			c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
			return
		}

		var isEmployeeInTransactions bool
		sqlCheckCustomerInTransaction := "SELECT EXISTS (SELECT 1 FROM bills WHERE employee_id = $1)"
		if err := db.QueryRow(sqlCheckCustomerInTransaction, id).Scan(&isEmployeeInTransactions); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check employee in transaction existence"})
			return
		}
		if isEmployeeInTransactions {
			c.JSON(http.StatusConflict, gin.H{"error": "Employee id is being used in transcations"})
			return
		}

		sqlDeleteEmployee := "DELETE FROM employees WHERE id = $1"
		_, err := db.Exec(sqlDeleteEmployee, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete employee"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Employee deleted successfully", "data": "OK"})
	}
}
