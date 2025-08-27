package model

type OrderItem struct {
	ID int64 `gorm:"primaryKey" json:"id"`
	OrderID int64 `json:"orderId"`
	ProductID int64 `json:"productId"`
	Product Product `gorm:"foreignKey:ProductID"`
	Qty int `json:"qty"`
	UnitPrice float64 `json:"unitPrice"`
}

