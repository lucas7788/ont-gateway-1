package model

// App for app
type App struct {
	ID               int    `bson:"id" json:"id"`
	Name             string `bson:"name" json:"name"`
	TxNotifyURL      string `bson:"tx_notify_url" json:"tx_notify_url"`
	PaymentNotifyURL string `bson:"payment_notify_url" json:"payment_notify_url"`
}
