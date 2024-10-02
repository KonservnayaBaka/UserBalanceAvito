package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"userBalanceAvito/internal/handlers"
)

func InitializeRoutes(r *gin.Engine, db *gorm.DB) {
	balanceGroup := r.Group("/balance")
	{
		balanceGroup.GET("/getBalance", handlers.GetBalance(db))
		balanceGroup.PUT("/enrollment", handlers.EnrollmentBalance(db))
		balanceGroup.PUT("/writeOff", handlers.WriteOffBalance(db))
		balanceGroup.PUT("/transfer", handlers.TransferBalance(db))
	}
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/registration", handlers.RegistrationNewUser(db))
		authGroup.GET("/login", handlers.AuthorizationUser(db))
	}
	productGroup := r.Group("/product")
	{
		productGroup.POST("/add", handlers.AddProduct(db))
		productGroup.POST("/reservation", handlers.ReservProduct(db))
	}
	operationsGroup := r.Group("/operations")
	{
		operationsGroup.POST("/create", handlers.CreateOperations(db))
		operationsGroup.GET("/get", handlers.GetOperations(db))
	}
}
