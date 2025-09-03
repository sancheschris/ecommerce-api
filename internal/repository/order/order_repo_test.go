package order

import (
	"testing"

	"github.com/sancheschris/ecommerce-api/internal/model"
	"github.com/sancheschris/ecommerce-api/internal/repository"
	"github.com/stretchr/testify/assert"
)


func TestCreateNewOrder(t *testing.T) {
	db := repository.SetupTestDB(model.Order{}, model.OrderItem{})

	orderDB := NewOrder(db)

	tests := []struct {
			name      string
			userID    int64
			items     []model.OrderItem
			status    string
			totalPrice float64
			currency  string
			payments  []model.Payment
			wantErr   bool
		}{
			{
				name:      "Valid order",
				userID:    1,
				items:     []model.OrderItem{{ProductID: 1, Qty: 2, UnitPrice: 100}},
				status:    "pending",
				totalPrice: 200,
				currency:  "USD",
				payments:  []model.Payment{},
				wantErr:   false,
			},
			{
				name:      "Missing user ID",
				userID:    0,
				items:     []model.OrderItem{{ProductID: 1, Qty: 2, UnitPrice: 100}},
				status:    "pending",
				totalPrice: 200,
				currency:  "USD",
				payments:  []model.Payment{},
				wantErr:   true,
			},
			{
				name:      "No items",
				userID:    1,
				items:     []model.OrderItem{},
				status:    "pending",
				totalPrice: 0,
				currency:  "USD",
				payments:  []model.Payment{},
				wantErr:   true,
			},
			{
				name:      "Negative quantity",
				userID:    1,
				items:     []model.OrderItem{{ProductID: 1, Qty: -2, UnitPrice: 100}},
				status:    "pending",
				totalPrice: -200,
				currency:  "USD",
				payments:  []model.Payment{},
				wantErr:   true,
			},
		}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
				order, err := model.NewOrder(tt.userID, tt.items, tt.status, tt.totalPrice, tt.currency, tt.payments)
			if tt.wantErr {
				assert.Error(t, err)
					return
				}
				assert.NoError(t, err)
				err = orderDB.CreateOrder(order)
				assert.NoError(t, err)
		})
	}
}
func TestGetOrders(t *testing.T) {
	db := repository.SetupTestDB(model.Order{}, model.OrderItem{}, model.Payment{})
	orderDB := NewOrder(db)

	// Seed some orders
	order1, _ := model.NewOrder(1, []model.OrderItem{{ProductID: 1, Qty: 2, UnitPrice: 100}}, "pending", 200, "USD", []model.Payment{})
	order2, _ := model.NewOrder(2, []model.OrderItem{{ProductID: 2, Qty: 1, UnitPrice: 200}}, "pending", 200, "USD", []model.Payment{})
	_ = orderDB.CreateOrder(order1)
	_ = orderDB.CreateOrder(order2)

	orders, err := orderDB.GetOrders()
	assert.NoError(t, err)
	assert.Len(t, orders, 2)
	
	for i, expected := range []model.Order{*order1, *order2} {
		actual := orders[i]
		assert.Equal(t, expected.ID, actual.ID)
		assert.Equal(t, expected.UserID, actual.UserID)
		assert.Equal(t, expected.Status, actual.Status)
		assert.Equal(t, expected.TotalPrice, actual.TotalPrice)
		assert.Equal(t, expected.Currency, actual.Currency)
		assert.Equal(t, expected.Items, actual.Items)
		assert.Equal(t, expected.Payments, actual.Payments)
		assert.True(t, expected.CreatedAt.Equal(actual.CreatedAt))
		assert.True(t, expected.UpdatedAt.Equal(actual.UpdatedAt))
	}
}
