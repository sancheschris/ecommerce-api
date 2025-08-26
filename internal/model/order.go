package model

import "time"

type Order struct {
	ID int64 `json:"id" gorm:"primaryKey"`
	UserID int64 `json:"userId"`
	User User `gorm:"foreignKey:UserID"`
	Status string `json:"status"`
	TotalPrice float64 `json:"total_price"`
	Currency string `json:"currency"`
	Items []OrderItem `gorm:"foreignKey:OrderID" json:"items"`
	Payments []Payment `gorm:"foreignKey:OrderID" json:"payments"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}