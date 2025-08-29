package model

import (
	"errors"
	"testing"
)

func TestNewProduct(t *testing.T) {
	tests := []struct {
		name string
		product Product
		wantErr error
	} {
		{"ok", Product{Name: "Iphone 16", Price: 7900.00, Active: true}, nil},
		{"missing name", Product{Name: "", Price: 7900.00, Active: true}, ErrNameIsRequired},
		{"missing price", Product{Name: "Iphone 16", Price: 0.0, Active: true}, ErrPriceIsRequired},
		{"invalid price", Product{Name: "Iphone 16", Price: -100.00, Active: true}, ErrInvalidPrice},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
				ans := tt.product.ValidateFields()
				if ans != tt.wantErr {
					t.Errorf("got %v, want %v", ans, tt.wantErr)
				}
				if !errors.Is(ans, tt.wantErr) {
					t.Fatalf("got error %v, want %v", ans, tt.wantErr)
				}
			})
	}
}