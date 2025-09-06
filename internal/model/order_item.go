package model

type OrderItem struct {
	ID int `gorm:"primaryKey" json:"id"`
	OrderID int `json:"orderId"`
	ProductID int `json:"productId"`
	Product Product `gorm:"foreignKey:ProductID"`
	Qty int `json:"qty"`
	UnitPrice float64 `json:"unitPrice"`
}

