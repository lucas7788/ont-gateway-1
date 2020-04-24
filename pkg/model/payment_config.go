package model

// PayPeriod for pay period
type PayPeriod int

const (
	// PayPeriodOnce for once payment
	PayPeriodOnce PayPeriod = iota + 1
	// PayPeriodDaily for daily payment
	PayPeriodDaily
	// PayPeriodMonthly for monthly payment
	PayPeriodMonthly
	// PayPeriodSeasonly for seasonly payment
	PayPeriodSeasonly
	// PayPeriodYearly for yearly payment
	PayPeriodYearly
)

// CoinType for coin type
type CoinType int

const (
	// CoinTypeONG for ONG
	CoinTypeONG CoinType = iota
)

// PaymentConfig for payment config
type PaymentConfig struct {
	App             int         `bson:"app" json:"app"`
	PaymentConfigID string      `bson:"payment_config_id" json:"payment_config_id"`
	AmountOptions   []int       `bson:"amount_options" json:"amount_options"`
	PeriodOptions   []PayPeriod `bson:"period_options" json:"period_options"`
	CoinType        CoinType    `bson:"coin_type" json:"coin_type"`
	PayMethods      []PayMethod `bson:"pay_methods" json:"pay_methods"`
	p2a             map[PayPeriod]int
}

// HasPayMethod checks whether PayMethod is allowed
func (config *PaymentConfig) HasPayMethod(pm PayMethod) bool {
	for _, payMethod := range config.PayMethods {
		if payMethod == pm {
			return true
		}
	}
	return false
}

// AmountForPeriod returns the amount for period
func (config *PaymentConfig) AmountForPeriod(period PayPeriod) (amount int, defined bool) {
	if config.p2a == nil {
		p2a := make(map[PayPeriod]int)
		for i, amount := range config.AmountOptions {
			period := config.PeriodOptions[i]
			p2a[period] = amount
		}
		config.p2a = p2a
	}

	amount, defined = config.p2a[period]
	return
}
