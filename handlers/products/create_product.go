package products

import (
	"challenge-goapi/models"
	"challenge-goapi/util"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Post(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newProduct models.Product

		if err := c.ShouldBindJSON(&newProduct); err != nil {
			errors := util.Validate(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": errors})
			return
		}

		sqlInsertProduct := "INSERT INTO products (name, price, unit) VALUES ($1, $2, $3) RETURNING id"
		if err := db.QueryRow(sqlInsertProduct, newProduct.Name, newProduct.Price, newProduct.Unit).Scan(&newProduct.Id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product in database."})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Product created successfully", "data": newProduct})
	}
}
