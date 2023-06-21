package entity

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	id           string
	sellingOrder *Order
	buyingOrder  *Order
	shares       int
	price        float64
	total        float64
	datetime     time.Time
}

func NewTransaction(sellingOrder *Order, buyingOrder *Order, shares int, price float64) *Transaction {
	total := float64(shares) * price
	return &Transaction{
		id:           uuid.New().String(),
		sellingOrder: sellingOrder,
		buyingOrder:  buyingOrder,
		shares:       shares,
		price:        price,
		total:        total,
		datetime:     time.Now(),
	}
}
