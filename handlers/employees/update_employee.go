package employees

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

		var updateEmployee models.EmployeUpdate
		if err := c.ShouldBindJSON(&updateEmployee); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
			return
		}

		var currentEmployee models.Employee
		sqlGetEmployee := "SELECT id, name, phone_number, address FROM employees WHERE id = $1;"
		if err := db.QueryRow(sqlGetEmployee, id).Scan(&currentEmployee.Id, &currentEmployee.Name, &currentEmployee.PhoneNumber, &currentEmployee.Address); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve current employee data"})
			return
		}

		if updateEmployee.Name != nil && strings.TrimSpace(*updateEmployee.Name) != "" {
			currentEmployee.Name = *updateEmployee.Name
		}
		if updateEmployee.PhoneNumber != nil && strings.TrimSpace(*updateEmployee.PhoneNumber) != "" {
			currentEmployee.PhoneNumber = *updateEmployee.PhoneNumber
		}
		if updateEmployee.Address != nil && strings.TrimSpace(*updateEmployee.Address) != "" {
			currentEmployee.Address = *updateEmployee.Address
		}

		sqlUpdateCustomer := "UPDATE employees SET name = $2, phone_number = $3, address = $4 WHERE id = $1"
		_, err := db.Exec(sqlUpdateCustomer, currentEmployee.Id, currentEmployee.Name, currentEmployee.PhoneNumber, currentEmployee.Address)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update employee"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Employee updated successfully", "data": currentEmployee})
	}
}

/*
Request :

- Method : PUT
- Endpoint : `/employees/:id` atau `/employees/` jika semua data
- Header :
  - Content-Type : application/json
  - Accept : application/json
- Body :

```json
{
  "name": "string",
  "phoneNumber": "string",
  "address": "string"
}
```

Response :

- Status : 200 OK
- Body :

```json
{
  "message": "string",
  "data": {
    "id": "string",
    "name": "string",
    "phoneNumber": "string",
    "address": "string"
  }
}
```
*/