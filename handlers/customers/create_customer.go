package customers

import (
	"challenge-goapi/models"
	"challenge-goapi/util"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Post(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newCustomer models.Customer

		if err := c.ShouldBindJSON(&newCustomer); err != nil {
			errors := util.Validate(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": errors})
			return
		}

		sqlInsertCustomer := "INSERT INTO customers (name, phone_number, address) VALUES ($1, $2, $3) RETURNING id;"
		if err := db.QueryRow(sqlInsertCustomer, newCustomer.Name, newCustomer.PhoneNumber, newCustomer.Address).Scan(&newCustomer.Id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create customer in database."})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Customer created successfully", "data": newCustomer})
	}
}
