package main

import (
	"database/sql"
	"log"

	"github.com/criotech/go-product-api/internal/controllers"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	router := gin.Default()
	address := ":3004"
	var db *sql.DB
	var e error

	if db, e = sql.Open("sqlite3", "./data.db"); e != nil {
		log.Fatalf("Error: %v", e)
	}
	defer db.Close()

	if e := db.Ping(); e != nil {
		log.Fatalf("Error: %v", e)
	}

	router.GET("/products", controllers.GetProducts(db))
	router.GET("/products/:guid", controllers.GetProduct(db))
	router.POST("/products", controllers.PostProduct(db))
	router.DELETE("/products/:guid", controllers.DeleteProduct(db))
	router.PUT("/products/:guid", controllers.PutProduct(db))

	log.Fatalln(router.Run((address)))
}
