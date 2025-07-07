package frontend

import (
	"context"
	_ "embed"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	//go:embed sql/get_balance.sql
	getBalanceQuery string
)

type Repository interface {
	GetBalance(context.Context)
}

type WalletRepository struct {
	pool *pgxpool.Pool
}

func NewWalletRepository(pool *pgxpool.Pool) *WalletRepository {
	return &WalletRepository{pool: pool}
}

func (r *WalletRepository) GetBalance(c context.Context, walletId string) (Wallet, error) {
	var wallet Wallet

	err := r.pool.QueryRow(c, getBalanceQuery, walletId).
		Scan(
			&wallet.Id,
			&wallet.Balances,
		)

	if err != nil {
		return Wallet{}, err
	}

	return wallet, nil
}
