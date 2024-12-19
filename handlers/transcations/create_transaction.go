package transcations

import (
	"challenge-goapi/models"
	"challenge-goapi/util"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Post(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newBill models.Bill
		if err := c.ShouldBindJSON(&newBill); err != nil {
			errors := util.Validate(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": errors})
			return
		}

		var isCustomerExists bool
		sqlCheckCustomer := "SELECT EXISTS (SELECT 1 FROM customers WHERE id = $1)"
		if err := db.QueryRow(sqlCheckCustomer, newBill.CustomerID).Scan(&isCustomerExists); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check customer existence"})
			return
		}
		if !isCustomerExists {
			c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
			return
		}

		var isEmployeeExists bool
		sqlCheckEmployee := "SELECT EXISTS (SELECT 1 FROM employees WHERE id = $1)"
		if err := db.QueryRow(sqlCheckEmployee, newBill.EmployeeID).Scan(&isEmployeeExists); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check employee existence"})
			return
		}
		if !isEmployeeExists {
			c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
			return
		}

		tx, err := db.Begin()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
			return
		}

		sqlInsertBill := "INSERT INTO bills (bill_date, entry_date, finish_date, employee_id, customer_id) VALUES ($1, $2, $3, $4, $5) RETURNING id"
		if err := tx.QueryRow(sqlInsertBill, newBill.BillDate, newBill.EntryDate, newBill.FinishDate, newBill.EmployeeID, newBill.CustomerID).Scan(&newBill.Id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create bill"})
			tx.Rollback()
			return
		}

		for i, detail := range newBill.BillDetails {
			var productPrice int
			sqlCheckPrice := "SELECT price FROM products WHERE id = $1"
			if err := tx.QueryRow(sqlCheckPrice, detail.ProductID).Scan(&productPrice); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Failed to fetch price for product ID %s", detail.ProductID)})
				tx.Rollback()
				return
			}

			sqlInsertBillDetails := "INSERT INTO bill_details (bill_id, product_id, product_price, qty) VALUES ($1, $2, $3, $4) RETURNING id"
			if err := tx.QueryRow(sqlInsertBillDetails, newBill.Id, detail.ProductID, productPrice, detail.Qty).Scan(&newBill.BillDetails[i].Id); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create bill details"})
				tx.Rollback()
				return
			}

			newBill.BillDetails[i].BillID = newBill.Id
			newBill.BillDetails[i].ProductPrice = productPrice
		}

		if err := tx.Commit(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Transaction created", "data": newBill})
	}
}