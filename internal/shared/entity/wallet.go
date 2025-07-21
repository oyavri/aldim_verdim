package entity

type Wallet struct {
	Id       string    `json:"id"`
	Balances []Balance `json:"balances"`
}

type Balance struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

type Event struct {
	AppId            string           `json:"app" validate:"required"`
	ActionType       string           `json:"type" validate:"required,oneof=BALANCE_INCREASE BALANCE_DECREASE"`
	Time             string           `json:"time" validate:"required"` // needs refactoring according to Kafka or db
	Meta             Meta             `json:"meta"`
	WalletId         string           `json:"wallet" validate:"required"`
	ActionAttributes ActionAttributes `json:"attributes" validate:"required"`
}

type Meta struct {
	UserId string `json:"user"`
}

type ActionAttributes struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency" validate:"oneof=TRY USD"`
}
