package transcations

import (
	"challenge-goapi/models"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Get(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		billID := c.Param("id_bill")
		if billID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or missing bill ID"})
			return
		}


		var bill models.Bill
		sqlGetBill := "SELECT id, bill_date, entry_date, finish_date, employee_id, customer_id FROM bills WHERE id = $1"
		if err := db.QueryRow(sqlGetBill, billID).Scan(&bill.Id, &bill.BillDate, &bill.EntryDate, &bill.FinishDate, &bill.EmployeeID, &bill.CustomerID); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
			return
		}

		var employee models.Employee
		sqlGetEmployee := "SELECT id, name, phone_number, address FROM employees WHERE id = $1"
		if err := db.QueryRow(sqlGetEmployee, bill.EmployeeID).Scan(&employee.Id, &employee.Name, &employee.PhoneNumber, &employee.Address); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch employee data"})
			return
		}

		var customer models.Customer
		sqlGetCustomer := "SELECT id, name, phone_number, address FROM customers WHERE id = $1"
		if err := db.QueryRow(sqlGetCustomer, bill.CustomerID).Scan(&customer.Id, &customer.Name, &customer.PhoneNumber, &customer.Address); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch customer data"})
			return
		}

		sqlGetBillDetailsWithProducts := "SELECT bill_details.id, bill_details.bill_id, bill_details.product_id, bill_details.product_price, bill_details.qty, products.id, products.name, products.price, products.unit FROM bill_details JOIN products ON bill_details.product_id = products.id WHERE bill_details.bill_id = $1;"
		rows, err := db.Query(sqlGetBillDetailsWithProducts, bill.Id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bill details"})
			return
		}
		defer rows.Close()

		var totalBill int
		var billDetails []gin.H
		for rows.Next() {
			var bd models.BillDetails
			var p models.Product

			if err = rows.Scan(&bd.Id, &bd.BillID, &bd.ProductID, &bd.ProductPrice, &bd.Qty, &p.Id, &p.Name, &p.Price, &p.Unit); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan bill detail"})
				return
			}

			totalBill += bd.ProductPrice * bd.Qty
			billDetails = append(billDetails, gin.H{
				"id":     bd.Id,
				"billId": bd.BillID,
				"product": gin.H{
					"id":    p.Id,
					"name":  p.Name,
					"price": p.Price,
					"unit":  p.Unit,
				},
				"productPrice": bd.ProductPrice,
				"qty":          bd.Qty,
			})
		}

		if err := rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to iterate bill details rows"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Transaction retrieved successfully",
			"data": gin.H{
				"id":         bill.Id,
				"billDate":   bill.BillDate,
				"entryDate":  bill.EntryDate,
				"finishDate": bill.FinishDate,
				"employee": gin.H{
					"id":          employee.Id,
					"name":        employee.Name,
					"phoneNumber": employee.PhoneNumber,
					"address":     employee.Address,
				},
				"customer": gin.H{
					"id":          customer.Id,
					"name":        customer.Name,
					"phoneNumber": customer.PhoneNumber,
					"address":     customer.Address,
				},
				"billDetails": billDetails,
				"totalBill":   totalBill,
			},
		})
	}
}