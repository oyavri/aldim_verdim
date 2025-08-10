package entity

import "time"

type Wallet struct {
	Id       string    `json:"id"`
	Balances []Balance `json:"balances"`
}

type Balance struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}

type Event struct {
	AppId            string           `json:"app"`
	ActionType       string           `json:"type"`
	Time             time.Time        `json:"time"` // needs refactoring according to Kafka or db
	Meta             Meta             `json:"meta"`
	WalletId         string           `json:"wallet"`
	ActionAttributes ActionAttributes `json:"attributes"`
}

type Meta struct {
	UserId string `json:"user"`
}

type ActionAttributes struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}
