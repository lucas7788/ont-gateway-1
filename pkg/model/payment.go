package model

import (
	"fmt"
	"time"
)

// PaymentState for payment state
type PaymentState int

const (
	// PaymentStateToStart when paid amount is not enough to get started
	PaymentStateToStart PaymentState = iota
	// PaymentStateStarted when paid amount is enough to get started
	PaymentStateStarted
	// PaymentStateEnd when payment ended
	PaymentStateEnd
)

// Payment model
type Payment struct {
	App               int          `bson:"app" json:"app"`
	PaymentID         string       `bson:"payment_id" json:"payment_id"`
	PaymentConfigID   string       `bson:"payment_config_id" json:"payment_config_id"`
	PayPeriod         PayPeriod    `bson:"pay_period" json:"pay_period"`
	PayMethod         PayMethod    `bson:"pay_method" json:"pay_method"`
	UnitAmount        int          `bson:"unit_amount" json:"unit_amount"`
	State             PaymentState `bson:"state" json:"state"`
	StartTime         time.Time    `bson:"start_time" json:"start_time"`
	EndTime           time.Time    `bson:"end_time" json:"end_time"`
	Balance           int          `bson:"balance" json:"balance"`
	BalanceExpireTime time.Time    `bson:"balance_expire_time" json:"balance_expire_time"`
	NotifyDate        string       `bson:"notify_date" json:"notify_date"`
	UpdatedAt         time.Time    `bson:"updated_at" json:"updated_at"`
	NotifyErrCount    int          `bson:"notify_err_count" json:"notify_err_count"`
	NotifyErrMsg      string       `bson:"notify_err_msg" json:"notify_err_msg"`
}

// PayMethod for pay method
type PayMethod int

const (
	// PayMethodBeforeUse for pay before use
	PayMethodBeforeUse PayMethod = iota
	// PayMethodAfterUse for pay after use
	PayMethodAfterUse
)

// IsStarted tells whether payment is started
func (p *Payment) IsStarted() bool {
	return p.State == PaymentStateStarted
}

// IsEnded tells whether payment is ended
func (p *Payment) IsEnded() bool {
	return p.State == PaymentStateEnd
}

// PeriodDuration returns the period duration
func (p *Payment) PeriodDuration() (d time.Duration) {
	switch p.PayPeriod {
	case PayPeriodDaily:
		d = time.Hour * 24
	case PayPeriodMonthly:
		d = time.Hour * 24 * 30
	case PayPeriodSeasonly:
		d = time.Hour * 24 * 30 * 3
	case PayPeriodYearly:
		d = time.Hour * 24 * 30 * 3 * 4
	case PayPeriodOnce:
		// 100 years
		d = time.Hour * 24 * 365 * 100
	}

	return
}

// PeriodNth returns current Period Nth
func (p *Payment) PeriodNth() (n int) {

	emptyTime := time.Time{}
	if p.StartTime == emptyTime {
		return
	}

	duration := time.Now().Sub(p.StartTime)

	switch p.PayPeriod {
	case PayPeriodDaily:
		n = int(duration / time.Hour / 24)
	case PayPeriodMonthly:
		n = int(duration / time.Hour / 24 / 30)
	case PayPeriodSeasonly:
		n = int(duration / time.Hour / 24 / 30 / 3)
	case PayPeriodYearly:
		n = int(duration / time.Hour / 24 / 30 / 3 / 4)
	}

	return
}

// AmountForNth returns the total amount for the n-th(zero based) PayPeriod
func (p *Payment) AmountForNth(n int, paymentConfig *PaymentConfig) (amount int, err error) {
	unitAmount, exists := paymentConfig.AmountForPeriod(p.PayPeriod)
	if !exists {
		err = fmt.Errorf("non existing PayPeriod %v for PaymentConfig %v |", p.PayPeriod, paymentConfig.PaymentConfigID)
		return
	}
	if p.PayPeriod == PayPeriodOnce {
		amount = unitAmount
		return
	}
	switch p.PayMethod {
	case PayMethodBeforeUse:
		amount = (n + 1) * unitAmount
	case PayMethodAfterUse:
		amount = n * unitAmount
	}
	return
}
