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

func (t *Transaction) CalculateTotal(shares int, price float64) {
	t.total = float64(shares) * t.price
}

func (t *Transaction) CloseBuyOrder() {
	if t.buyingOrder.pedingShares <= 0 {
		t.buyingOrder.status = "CLOSED"
	}
}

func (t *Transaction) CloseSellOrder() {
	if t.sellingOrder.pedingShares == 0 {
		t.sellingOrder.status = "CLOSED"
	}
}

func (t *Transaction) AddBuyOrderPendingShares(shares int) {
	t.buyingOrder.pedingShares += shares
}

func (t *Transaction) AddSellOrderPendingShares(shares int) {
	t.buyingOrder.pedingShares += shares
}
