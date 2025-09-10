package model

type OrderItem struct {
    ID         int       `gorm:"primaryKey" json:"id"`
    OrderID    int       `json:"order_id"`
    ProductID  int       `json:"product_id"`
    Product    *Product  `json:"product,omitempty"` 
    Qty        int       `json:"qty"`
    UnitPrice  float64   `json:"unit_price"`
}