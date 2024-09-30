package models

import "time"

type User struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

type Account struct {
	ID        uint      `json:"id"`
	Balance   int       `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

type Product struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type Reservation struct {
	ID        uint      `json:"id"`
	AccountID uint      `json:"account_id"`
	ProductID uint      `json:"product_id"`
	OrderID   uint      `json:"order_id"`
	Amount    int       `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
	Account   Account   `json:"account"`
	Product   Product   `json:"product"`
}

type Operation struct {
	ID            uint      `json:"id"`
	AccountID     uint      `json:"account_id"`
	Amount        int       `json:"amount"`
	OperationType string    `json:"operation_type"`
	CreatedAt     time.Time `json:"created_at"`
	ProductID     *uint     `json:"product_id,omitempty"`
	OrderID       *uint     `json:"order_id,omitempty"`
	Description   *string   `json:"description,omitempty"`
	Account       Account   `json:"account"`
}
