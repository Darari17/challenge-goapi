package main

import (
	"challenge-goapi/config"
	"challenge-goapi/handlers/customers"
	"challenge-goapi/handlers/employees"
	"challenge-goapi/handlers/products"
	"challenge-goapi/handlers/transcations"

	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	db := config.ConnectDB()
	if db == nil {
		fmt.Println("Failed to establish database connection. Exiting program.")
		return
	}
	defer db.Close()

	router := gin.Default()

	customerGroup := router.Group("/customers")
	{
		customerGroup.POST("/", customers.Post(db))
		customerGroup.GET("/", customers.Get(db))
		customerGroup.GET("/:id", customers.Get(db))
		customerGroup.PUT("/:id", customers.Put(db))
		customerGroup.DELETE("/:id", customers.Delete(db))
	}

	employeeGroup := router.Group("/employees")
	{
		employeeGroup.POST("/", employees.Post(db))
		employeeGroup.GET("/", employees.Get(db))
		employeeGroup.GET("/:id", employees.Get(db))
		employeeGroup.PUT("/:id", employees.Put(db))
		employeeGroup.DELETE("/:id", employees.Delete(db))

	}

	productGroup := router.Group("/products")
	{
		productGroup.POST("/", products.Post(db))
		productGroup.GET("/", products.Get(db))
		productGroup.GET("/:id", products.Get(db))
		productGroup.PUT("/:id", products.Put(db))
		productGroup.DELETE("/:id", products.Delete(db))
	}

	transcationsGroup := router.Group("/transactions")
	{
		transcationsGroup.POST("/", transcations.Post(db))
		transcationsGroup.GET("/:id_bill", transcations.Get(db))
		transcationsGroup.GET("/", transcations.ListTransactions(db))
	}

	router.Run(":8080")
}