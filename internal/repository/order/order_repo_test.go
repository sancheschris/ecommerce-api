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

func TestGetOrdersByUserID(t *testing.T) {
	db := repository.SetupTestDB(model.Order{}, model.OrderItem{})
	orderDB := NewOrder(db)

	userID := int64(1)

	order, err := model.NewOrder(userID, []model.OrderItem{{ProductID: 1, Qty: 2, UnitPrice:  100},
	},
	"pending",
	200,
	"USD",
	[]model.Payment{},)

	assert.NoError(t, err)
	assert.NotNil(t, order)

	orderDB.CreateOrder(order)

	orderByID, err := orderDB.GetOrdersByUserID(order.UserID)

	assert.NoError(t, err)
	assert.NotNil(t, orderByID)
	assert.Equal(t, int64(1), orderByID[0].UserID)
	assert.Equal(t, order.Status, orderByID[0].Status)

}

func TestGetOrderByID(t *testing.T) {
	db := repository.SetupTestDB(model.Order{}, model.OrderItem{})
	orderDB := NewOrder(db)

	order, _ := model.NewOrder(1, []model.OrderItem{{ProductID: 1, Qty: 2, UnitPrice: 100}}, "pending", 200, "USD", []model.Payment{})

	orderDB.CreateOrder(order)

	actual, err := orderDB.GetOrderByID(order.ID)

	if err != nil {
		t.Error(err)
	}

	if actual.TotalPrice != order.TotalPrice {
		t.Errorf("actual totalPrice: %f, want 200", actual.TotalPrice)
	}

}

func TestUpdateOrder(t *testing.T) {
	db := repository.SetupTestDB(model.Order{}, model.OrderItem{})
	orderDB := NewOrder(db)

	order, _ := model.NewOrder(1, []model.OrderItem{{ProductID: 1, Qty: 2, UnitPrice: 100}}, "pending", 200, "USD", []model.Payment{})
	orderDB.CreateOrder(order)

	existingOrder, _ := orderDB.GetOrderByID(order.ID)
	assert.Equal(t, 200.00, existingOrder.TotalPrice)
	assert.Equal(t, "pending", existingOrder.Status)

	existingOrder.TotalPrice = 350.00
	existingOrder.Status = "Done"

	orderDB.UpdateOrder(existingOrder)

	actual, _ := orderDB.GetOrderByID(existingOrder.ID)

	assert.NotNil(t, actual)
	assert.Equal(t, "Done", actual.Status)
	assert.Equal(t, 350.00, actual.TotalPrice)
	assert.Equal(t, "USD", actual.Currency)

}

func TestDeleteOrder(t *testing.T) {
	db := repository.SetupTestDB(model.Order{}, model.OrderItem{})
	orderDB := NewOrder(db)

	order, _ := model.NewOrder(1, []model.OrderItem{{ProductID: 1, Qty: 2, UnitPrice: 100}}, "pending", 200, "USD", []model.Payment{})
	orderDB.CreateOrder(order)

	existingOrder, _ := orderDB.GetOrderByID(order.ID)

	assert.Equal(t, "USD", existingOrder.Currency)

	err := orderDB.DeleteOrder(existingOrder.ID)
	assert.NoError(t, err)

	deletedOrder, _ := orderDB.GetOrderByID(existingOrder.ID)
	assert.Nil(t, deletedOrder)

}

func TestAddOrderItem(t *testing.T) {
	db := repository.SetupTestDB(model.Order{}, model.OrderItem{})
	orderDB := NewOrder(db)

	order, _ := model.NewOrder(1, []model.OrderItem{{ProductID: 1, Qty: 2, UnitPrice: 100}}, "pending", 200, "USD", []model.Payment{})
	orderDB.CreateOrder(order)

	newItem := model.OrderItem{ProductID: 2, Qty: 1, UnitPrice: 50}
	err := orderDB.AddOrderItem(order.ID, &newItem)
	assert.NoError(t, err)

	updatedOrder, _ := orderDB.GetOrderByID(order.ID)
	assert.Len(t, updatedOrder.Items, 2)
	assert.Equal(t, int64(2), updatedOrder.Items[1].ProductID)
	assert.Equal(t, 50.0, updatedOrder.Items[1].UnitPrice)
}

func TestUpdateOrderItem(t *testing.T) {
	db := repository.SetupTestDB(model.Order{}, model.OrderItem{})
	orderDB := NewOrder(db)

	order, _ := model.NewOrder(1, []model.OrderItem{{ProductID: 1, Qty: 2, UnitPrice: 100}}, "pending", 200, "USD", []model.Payment{})
	orderDB.CreateOrder(order)

	newItem := model.OrderItem{ProductID: 2, Qty: 1, UnitPrice: 50.0}
	err := orderDB.AddOrderItem(order.ID, &newItem)
	assert.NoError(t, err)
	assert.Equal(t, 50.0, newItem.UnitPrice)

	newItem.UnitPrice = 35.0
	newItem.Qty = 5

	err = orderDB.UpdateOrderItem(order.ID, &newItem)
	assert.NoError(t, err)

	actual, err := orderDB.GetOrderItems(order.ID)
	assert.NoError(t, err)

	assert.Equal(t, 35.0, actual[1].UnitPrice)
	assert.Equal(t, 5, actual[1].Qty)
	assert.Equal(t, int64(2), actual[1].ProductID)
}

func TestGetOrderItems(t *testing.T) {
	db := repository.SetupTestDB(model.Order{}, model.OrderItem{})
	orderDB := NewOrder(db)

	order, _ := model.NewOrder(1, []model.OrderItem{{ProductID: 1, Qty: 2, UnitPrice: 100}}, "pending", 200, "USD", []model.Payment{})
	orderDB.CreateOrder(order)

	newItem := model.OrderItem{ProductID: 2, Qty: 1, UnitPrice: 50.0}
	err := orderDB.AddOrderItem(order.ID, &newItem)
	assert.NoError(t, err)

	orderItems, err := orderDB.GetOrderItems(order.ID)
	assert.NoError(t, err)
	assert.Equal(t, 100.00, orderItems[0].UnitPrice)
	assert.Equal(t, int64(1), orderItems[0].ProductID)
	assert.Equal(t, 50.00, orderItems[1].UnitPrice)
	assert.Equal(t, int64(2), orderItems[1].ProductID)
}

func TestRemoveOrderItem(t *testing.T) {
	db := repository.SetupTestDB(model.Order{}, model.OrderItem{})
	orderDB := NewOrder(db)

	order, _ := model.NewOrder(1, []model.OrderItem{{ProductID: 1, Qty: 2, UnitPrice: 100}}, "pending", 200, "USD", []model.Payment{})
	orderDB.CreateOrder(order)

	newItem := model.OrderItem{ProductID: 2, Qty: 1, UnitPrice: 50.0}
	err := orderDB.AddOrderItem(order.ID, &newItem)
	assert.NoError(t, err)

	// Remove the item
	err = orderDB.RemoveOrderItem(order.ID, newItem.ID)
	assert.NoError(t, err)

	// Verify the item is removed
	items, err := orderDB.GetOrderItems(order.ID)
	assert.NoError(t, err)
	assert.Len(t, items, 1)
	assert.Equal(t, int64(1), items[0].ProductID)

	// Try removing a non-existent item
	err = orderDB.RemoveOrderItem(order.ID, 9999)
	assert.Error(t, err)
}