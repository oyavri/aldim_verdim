package dto

import "github.com/oyavri/aldim_verdim/pkg/entity"

type EventRequest struct {
	Events []entity.Event `json:"events" validate:"required"`
}

type WalletResponse struct {
	Wallets []entity.Wallet `json:"wallets"`
}
