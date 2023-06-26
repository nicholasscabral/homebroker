package entity

import (
	"container/heap"
	"sync"
)

type Book struct {
	orders        []*Order
	transactions  []*Transaction
	ordersChan    chan *Order
	ordersChanOut chan *Order
	wg            *sync.WaitGroup
}

func NewBook(orderChan chan *Order, orderChanOut chan *Order, wg *sync.WaitGroup) *Book {
	return &Book{
		orders:        []*Order{},
		transactions:  []*Transaction{},
		ordersChan:    orderChan,
		ordersChanOut: orderChanOut,
		wg:            wg,
	}
}

func (b *Book) Trade() {
	buyOrders := NewOrderQueue()
	sellOrders := NewOrderQueue()

	heap.Init(buyOrders)
	heap.Init(sellOrders)

	for order := range b.ordersChan {
		if order.orderType == "buy" {
			buyOrders.Push(order)
			if sellOrders.Len() > 0 && sellOrders.Orders[0].price <= order.price {
				sellOrder := sellOrders.Pop().(*Order)
				if sellOrder.pedingShares > 0 {
					transaction := NewTransaction(sellOrder, order, order.shares, sellOrder.price)
					b.AddTransaction(transaction, b.wg)
					sellOrder.transaction = append(sellOrder.transaction, transaction)
					order.transaction = append(order.transaction, transaction)
					b.ordersChanOut <- sellOrder
					b.ordersChanOut <- order
					if sellOrder.pedingShares > 0 {
						sellOrders.Push(sellOrder)
					}
				}
			}
		} else if order.orderType == "sell" {
			sellOrders.Push(order)
			if buyOrders.Len() > 0 && buyOrders.Orders[0].price >= order.price {
				buyOrder := buyOrders.Pop().(*Order)
				if buyOrder.pedingShares > 0 {
					transaction := NewTransaction(order, buyOrder, order.shares, buyOrder.price)
					b.AddTransaction(transaction, b.wg)
					buyOrder.transaction = append(buyOrder.transaction, transaction)
					order.transaction = append(order.transaction, transaction)
					b.ordersChanOut <- buyOrder
					b.ordersChanOut <- order
					if buyOrder.pedingShares > 0 {
						buyOrders.Push(buyOrder)
					}
				}
			}
		}
	}
}

func (b *Book) AddTransaction(transaction *Transaction, wg *sync.WaitGroup) {
	defer wg.Done()

	sellingShares := transaction.sellingOrder.pedingShares
	buyingShares := transaction.buyingOrder.pedingShares

	minShares := sellingShares
	if buyingShares < sellingShares {
		minShares = buyingShares
	}

	transaction.sellingOrder.investor.UpdateAssetPosition(transaction.sellingOrder.asset.id, -minShares)
	transaction.AddSellOrderPendingShares(-minShares)

	transaction.buyingOrder.investor.UpdateAssetPosition(transaction.buyingOrder.asset.id, minShares)
	transaction.AddBuyOrderPendingShares(-minShares)

	transaction.CalculateTotal(transaction.shares, transaction.buyingOrder.price)
	transaction.CloseBuyOrder()
	transaction.CloseSellOrder()
	b.transactions = append(b.transactions, transaction)
}
