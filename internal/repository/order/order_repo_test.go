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
		name string
		order *model.Order
		wantErr bool
	} {
		{name: "Valid order",
		 order: &model.Order{UserID: 1, Items: []model.OrderItem{{ProductID: 1, Qty: 2, UnitPrice: 100}}},
		 wantErr: false,
		},
		{name: "Missing user ID",
		 order: &model.Order{Items: []model.OrderItem{{ProductID: 1, Qty: 2, UnitPrice: 100}}},
		 wantErr: true,
		},
		{name: "No items",
		 order: &model.Order{UserID: 1, Items: []model.OrderItem{}},
		 wantErr: true,
		},
		{
			name: "Negative quantity",
			order: &model.Order{UserID: 1, Items: []model.OrderItem{{ProductID: 1, Qty: -2, UnitPrice: 100}}},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			order, err := model.NewOrder(tt.order.UserID, tt.order.Items)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				err = orderDB.CreateOrder(order)
				assert.NoError(t, err)
			}
		})
	}
}