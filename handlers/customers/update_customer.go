package customers

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

		var updatedCustomer models.CustomerUpdate
		if err := c.ShouldBindJSON(&updatedCustomer); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
			return
		}

		var currentCustomer models.Customer
		sqlGetCustomer := "SELECT id, name, phone_number, address FROM customers WHERE id = $1;"
		if err := db.QueryRow(sqlGetCustomer, id).Scan(&currentCustomer.Id, &currentCustomer.Name, &currentCustomer.PhoneNumber, &currentCustomer.Address,); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve current customer data"})
			return
		}

		if updatedCustomer.Name != nil && strings.TrimSpace(*updatedCustomer.Name) != "" {
			currentCustomer.Name = *updatedCustomer.Name
		}
		if updatedCustomer.PhoneNumber != nil && strings.TrimSpace(*updatedCustomer.PhoneNumber) != "" {
			currentCustomer.PhoneNumber = *updatedCustomer.PhoneNumber
		}
		if updatedCustomer.Address != nil && strings.TrimSpace(*updatedCustomer.Address) != "" {
			currentCustomer.Address = *updatedCustomer.Address
		}

		sqlUpdateCustomer := "UPDATE customers SET name = $2, phone_number = $3, address = $4 WHERE id = $1"
		_, err := db.Exec(sqlUpdateCustomer, currentCustomer.Id, currentCustomer.Name, currentCustomer.PhoneNumber, currentCustomer.Address);
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update customer"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Customer updated successfully", "data": currentCustomer})
	}
}
