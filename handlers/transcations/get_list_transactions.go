package transcations

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func readSqlQuery(filePath string) (string, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file: %s\n", err.Error())
		return "", err
	}
	return string(file), nil
}

func ListTransactions(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		startDate := c.Query("startDate")
		endDate := c.Query("endDate")
		productName := c.Query("productName")
		dateLayout := "02-01-2006"

		query := `
			SELECT
				bills.id,
				TO_CHAR(bills.bill_date, 'DD-MM-YYYY'),
				TO_CHAR(bills.entry_date, 'DD-MM-YYYY'),
				TO_CHAR(bills.finish_date, 'DD-MM-YYYY'),
				employees.id,
				employees.name,
				employees.phone_number,
				employees.address,
				customers.id,
				customers.name,
				customers.phone_number,
				customers.address,
				bill_details.id,
				bill_details.bill_id,
				products.id,
				products.name,
				products.price,
				products.unit,
				bill_details.product_price,
				bill_details.qty,
				COALESCE(SUM(bill_details.product_price * bill_details.qty) OVER (PARTITION BY bills.id), 0)
			FROM
				bills
			LEFT JOIN
				employees ON bills.employee_id = employees.id
			LEFT JOIN
				customers ON bills.customer_id = customers.id
			LEFT JOIN
				bill_details ON bills.id = bill_details.bill_id
			LEFT JOIN
				products ON bill_details.product_id = products.id
			WHERE
				1=1`

		var params []interface{}
		paramIndex := 1

		if startDate != "" {
			query += ` AND bills.bill_date >= $` + strconv.Itoa(paramIndex)
			startDateParsed, err := time.Parse(dateLayout, startDate)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid startDate format. Use dd-MM-yyyy"})
				return
			}
			params = append(params, startDateParsed)
			paramIndex++
		}

		if endDate != "" {
			query += ` AND bills.bill_date <= $` + strconv.Itoa(paramIndex)
			endDateParsed, err := time.Parse(dateLayout, endDate)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid endDate format. Use dd-MM-yyyy"})
				return
			}
			params = append(params, endDateParsed)
			paramIndex++
		}

		if productName != "" {
			query += ` AND products.name ILIKE $` + strconv.Itoa(paramIndex)
			params = append(params, "%"+productName+"%")
			paramIndex++
		}

		query += ` ORDER BY bills.id, bill_details.id`

		rows, err := db.Query(query, params...)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transactions"})
			return
		}
		defer rows.Close()

		transactionsMap := make(map[string]gin.H)

		for rows.Next() {
			var (
				billID, billDate, entryDate, finishDate                      string
				employeeID, employeeName, employeePhone, employeeAddress     string
				customerID, customerName, customerPhone, customerAddress     string
				billDetailID, productID, productName, productUnit            string
				productPrice, detailProductPrice, detailQty, totalBill       int
			)

			err := rows.Scan(
				&billID, &billDate, &entryDate, &finishDate,
				&employeeID, &employeeName, &employeePhone, &employeeAddress,
				&customerID, &customerName, &customerPhone, &customerAddress,
				&billDetailID, &billID, &productID, &productName, &productPrice, &productUnit,
				&detailProductPrice, &detailQty, &totalBill,
			)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan transactions"})
				return
			}

			if _, exists := transactionsMap[billID]; !exists {
				transactionsMap[billID] = gin.H{
					"id":         billID,
					"billDate":   billDate,
					"entryDate":  entryDate,
					"finishDate": finishDate,
					"employee": gin.H{
						"id":          employeeID,
						"name":        employeeName,
						"phoneNumber": employeePhone,
						"address":     employeeAddress,
					},
					"customer": gin.H{
						"id":          customerID,
						"name":        customerName,
						"phoneNumber": customerPhone,
						"address":     customerAddress,
					},
					"billDetails": []gin.H{},
					"totalBill":   totalBill,
				}
			}

			transaction := transactionsMap[billID]
			transaction["billDetails"] = append(transaction["billDetails"].([]gin.H), gin.H{
				"id":      billDetailID,
				"billId":  billID,
				"product": gin.H{
					"id":    productID,
					"name":  productName,
					"price": productPrice,
					"unit":  productUnit,
				},
				"productPrice": detailProductPrice,
				"qty":          detailQty,
			})
			transactionsMap[billID] = transaction
		}

		transactions := make([]gin.H, 0, len(transactionsMap))
		for _, transaction := range transactionsMap {
			transactions = append(transactions, transaction)
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Transactions retrieved successfully",
			"data":    transactions,
		})
	}
}
