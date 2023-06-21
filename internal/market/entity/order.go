package entity

type OrderType = string

const (
	buy  OrderType = "buy"
	sell OrderType = "sell"
)

type Order struct {
	id           string
	investor     *Investor
	asset        *Asset
	shares       int
	pedingShares int
	price        float64
	orderType    OrderType
	status       string
	transaction  []*Transaction
}

func NewOrder(id string, investor *Investor, asset *Asset, shares int, price float64, orderType OrderType) *Order {
	return &Order{
		id:        id,
		investor:  investor,
		asset:     asset,
		shares:    shares,
		price:     price,
		orderType: orderType,
	}
}
