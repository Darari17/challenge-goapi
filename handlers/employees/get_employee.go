package employees

import (
	"challenge-goapi/models"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Get(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var employee models.Employee

		if id != "" {
			sqlGetEmployeeById := "SELECT id, name, phone_number, address FROM employees WHERE id = $1"
			if err := db.QueryRow(sqlGetEmployeeById, id).Scan(&employee.Id, &employee.Name, &employee.PhoneNumber, &employee.Address); err != nil {
				if err == sql.ErrNoRows {
					c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
				} else {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
				}
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": "Get employee successfully", "data": employee})
			return
		}

		sqlSelectEmployee := "SELECT id, name, phone_number, address FROM employees"
		rows, err := db.Query(sqlSelectEmployee)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
		defer rows.Close()

		var employees []models.Employee
		for rows.Next() {
			if err := rows.Scan(&employee.Id, &employee.Name, &employee.PhoneNumber, &employee.Address); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning employee data"})
				return
			}
			employees = append(employees, employee)
		}

		if err := rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading employee data"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Get all employee successfully", "data": employees})
	}
}
