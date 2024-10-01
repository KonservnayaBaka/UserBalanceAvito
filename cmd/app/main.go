package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"userBalanceAvito/internal/models"
	"userBalanceAvito/internal/routes"
)

func main() {
	r := gin.Default()

	dsn := "host=localhost user=postgres password=90814263 dbname=userBalanceAvito port=1488 sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	err = db.AutoMigrate(&models.User{}, &models.Account{}, &models.Product{}, &models.Reservation{}, &models.Operation{})
	if err != nil {
		panic("failed to migrate database")
	}

	routes.InitializeRoutes(r, db)

	r.Run("localhost:8000")
}
