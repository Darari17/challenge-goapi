package employees

import (
	"challenge-goapi/models"
	"challenge-goapi/util"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Post(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newEmployee models.Employee

		if err := c.ShouldBindJSON(&newEmployee); err != nil {
			errors := util.Validate(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": errors})
			return
		}

		sqlInsertEmployee := "INSERT INTO employees (name, phone_number, address) VALUES ($1, $2, $3) RETURNING id"
		if err := db.QueryRow(sqlInsertEmployee, newEmployee.Name, newEmployee.PhoneNumber, newEmployee.Address).Scan(&newEmployee.Id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create employee in database."})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Employee created successfully", "data": newEmployee})
	}
}