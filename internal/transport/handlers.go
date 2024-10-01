package transport

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
	"userBalanceAvito/internal/models"
)

func InitializeRoutes(r *gin.Engine, db *gorm.DB) {
	balanceGroup := r.Group("/balance")
	{
		balanceGroup.GET("/getBalance", getBalance(db))
		balanceGroup.PUT("/enrollment", EnrollmentBalance(db))
		balanceGroup.PUT("/writeOff", writeOffBalance(db))
		balanceGroup.PUT("/transfer", transferBalance(db))
	}
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/registration", registrationNewUser(db))
		authGroup.GET("/authorization", authorizationUser(db))
	}
	productGroup := r.Group("/product")
	{
		productGroup.POST("/add", addProduct(db))
		productGroup.POST("/reservation", reservProduct(db))
	}
}

func reservProduct(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetHeader("UserID")
		productID := c.GetHeader("ProductID")
		ammountCount := c.GetHeader("Ammount")

		var reservation models.Reservation

		userIDParse, err := strconv.ParseInt(userID, 10, 64)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
		}
		productIDParse, err := strconv.ParseInt(productID, 10, 64)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
		}
		ammountCountParse, err := strconv.ParseInt(ammountCount, 10, 64)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
		}

		reservation.AccountID = uint(userIDParse)
		reservation.ProductID = uint(productIDParse)
		reservation.Amount = int(ammountCountParse)

		c.JSON(200, gin.H{"reservation": reservation})
	}
}

func addProduct(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var product models.Product

		if err := c.ShouldBind(&product); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		if product.Name == "" {
			c.JSON(400, gin.H{"error": "product name is empty"})
			return
		}

		if err := db.Create(&product).Error; err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"product": product})
	}
}

func registrationNewUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		if user.Username == "" || user.Password == "" {
			c.JSON(400, gin.H{"error": "username or password is empty"})
			return
		}

		if err := db.Create(&user).Error; err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		account := models.Account{}

		if err := db.Create(&account).Error; err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"user": user, "account": account})
	}
}

func authorizationUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		if user.Username == "" || user.Password == "" {
			c.JSON(400, gin.H{"error": "username or password is empty"})
			return
		}

		if err := db.Where("username = ? AND password = ?", user.Username, user.Password).First(&user).Error; err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"user": user, "status": true})
	}
}

func getBalance(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetHeader("UserID")

		if userID == "" {
			c.JSON(400, gin.H{"error": "UserID is empty"})
			return
		}

		account := models.Account{}
		if err := db.Where("id = ?", userID).First(&account).Error; err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"balance": account.Balance})
	}
}

func EnrollmentBalance(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetHeader("UserID")
		changeNumberStr := c.GetHeader("ChangeNumber")

		if userID == "" || changeNumberStr == "" {
			c.JSON(400, gin.H{"error": "UserID is empty"})
			return
		}

		parseСhangeNumber, err := strconv.ParseInt(changeNumberStr, 10, 64)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		account := models.Account{}
		if err := db.Where("id = ?", userID).First(&account).Error; err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		account.Balance += int(parseСhangeNumber)
		if err := db.Save(&account).Error; err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
		}

		c.JSON(200, gin.H{"balance": account.Balance})
	}
}

func writeOffBalance(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetHeader("UserID")
		changeNumberStr := c.GetHeader("ChangeNumber")

		if userID == "" || changeNumberStr == "" {
			c.JSON(400, gin.H{"error": "UserID is empty"})
			return
		}

		parseСhangeNumber, err := strconv.ParseInt(changeNumberStr, 10, 64)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		account := models.Account{}
		if err := db.Where("id = ?", userID).First(&account).Error; err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		if account.Balance > 0 {
			account.Balance -= int(parseСhangeNumber)
		} else {
			c.JSON(400, gin.H{"error": "balance is zero"})
			return
		}
		if err := db.Save(&account).Error; err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
		}

		c.JSON(200, gin.H{"balance": account.Balance})
	}
}

func transferBalance(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		recipientID := c.GetHeader("RecipientID")
		senderID := c.GetHeader("SenderID")
		changeNumberStr := c.GetHeader("ChangeNumber")

		if recipientID == "" || senderID == "" || changeNumberStr == "" {
			c.JSON(400, gin.H{"error": "RecipientID or SenderID or ChangeNumber is empty"})
			return
		}

		parseChangeNumber, err := strconv.ParseInt(changeNumberStr, 10, 64)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		accountRecipient := models.Account{}
		accountSender := models.Account{}

		if err := db.Where("id = ?", recipientID).First(&accountRecipient).Error; err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		if err := db.Where("id = ?", senderID).First(&accountSender).Error; err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		if accountSender.ID == accountRecipient.ID {
			c.JSON(400, gin.H{"error": "sender is the same as recipient"})
			return
		}

		accountRecipient.Balance += int(parseChangeNumber)
		if err := db.Save(&accountRecipient).Error; err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		if accountSender.Balance > 0 {
			accountSender.Balance -= int(parseChangeNumber)
		} else {
			c.JSON(400, gin.H{"error": "balance is zero"})
			return
		}
		if err := db.Save(&accountSender).Error; err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"SenderBalance": accountSender.Balance, "RecipientBalance": accountRecipient.Balance})
	}
}
