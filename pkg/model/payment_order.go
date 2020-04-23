package model

import "fmt"

// PaymentOrder is order for payment
type PaymentOrder struct {
	App       int      `bson:"app" json:"app"`
	PaymentID string   `bson:"payment_id" json:"payment_id"`
	OrderID   string   `bson:"order_id" json:"order_id"`
	Amount    int      `bson:"amount" json:"amount"`
	CoinType  CoinType `bson:"coin_type" json:"coin_type"`
	OrderInfo string   `bson:"order_info" json:"order_info"`
}

// VerifyOrderInfo verifies OrderInfo
func VerifyOrderInfo(amount int, coinType CoinType, OrderInfo string) (err error) {
	switch coinType {
	case CoinTypeONG:
		return
	default:
		err = fmt.Errorf("CoinType %v not supported yet", coinType)
		return
	}
}
