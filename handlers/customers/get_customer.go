package customers

import (
	"challenge-goapi/models"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Get(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var customer models.Customer

		if id != "" {
			sqlGetCustomerByID := "SELECT id, name, phone_number, address FROM customers WHERE id = $1"
			if err := db.QueryRow(sqlGetCustomerByID, id).Scan(&customer.Id, &customer.Name, &customer.PhoneNumber, &customer.Address); err != nil {
				if err == sql.ErrNoRows {
					c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
				} else {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
				}
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": "Get customer successfully", "data": customer})
			return
		}

		sqlGetAllCustomers := "SELECT id, name, phone_number, address FROM customers;"
		rows, err := db.Query(sqlGetAllCustomers)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
		defer rows.Close()

		var customers []models.Customer
		for rows.Next() {
			if err := rows.Scan(&customer.Id, &customer.Name, &customer.PhoneNumber, &customer.Address); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning customer data"})
				return
			}
			customers = append(customers, customer)
		}

		if err := rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading customer data"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Get all customer successfully", "data": customers})
	}
}
