package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/criotech/go-product-api/internal"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type postProduct struct {
	Name        string `json:"name" binding:"required"`
	Price       string `json:"price" binding:"required,gt=0"`
	Description string `json:"description" binding:"omitempty,max=250"`
}
type Product struct {
	GUID        string `json:"guid"`
	Name        string `json:"name"`
	Price       string `json:"price"`
	Description string `json:"description"`
	CreatedAt   string `json:"createdAt"`
}

func PostProduct(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload postProduct
		ctx := c.Request.Context()
		fmt.Println("Hello world", ctx)

		if e := c.ShouldBindJSON(&payload); e != nil {
			res := internal.NewHTTPResponse(http.StatusBadRequest, e)
			c.JSON(http.StatusBadRequest, res)
			return
		}

		guid := uuid.New().String()
		createdAt := time.Now().Format(time.RFC3339)
		if _, e := db.ExecContext(ctx, "INSERT INTO products(guid, name, price, description, createdAt) VALUES(?,?,?,?,?)", guid, payload.Name, payload.Price, payload.Description, createdAt); e != nil {
			res := internal.NewHTTPResponse(http.StatusBadRequest, e)
			c.JSON(http.StatusInternalServerError, res)
			return
		}

		var product Product
		row := db.QueryRow("SELECT guid,name,price,description,createdAt FROM products WHERE guid=?", guid)

		if e := row.Scan(&product.GUID, &product.Name, &product.Price, &product.Description, &product.CreatedAt); e != nil {
			res := internal.NewHTTPResponse(http.StatusBadRequest, e)

			c.JSON(http.StatusInternalServerError, res)
			return
		}

		res := internal.NewHTTPResponse(http.StatusCreated, product)

		c.Writer.Header().Add("Location", fmt.Sprintf("/products/%s", guid))
		c.JSON(http.StatusCreated, res)
		fmt.Println(payload)
	}
}
