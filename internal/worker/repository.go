package worker

import (
	"context"

	_ "embed"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	//go:embed sql/increase_balance.sql
	increaseBalanceQuery string
	//go:embed sql/decrease_balance.sql
	decreaseBalanceQuery string
)

type WalletRepository interface {
	IncreaseBalance(ctx context.Context, userId string, walletId string, amount float64, currency string) error
	DecreaseBalance(ctx context.Context, userId string, walletId string, amount float64, currency string) error
}

type WalletRepo struct {
	pool *pgxpool.Pool
}

func NewWalletRepository(pool *pgxpool.Pool) *WalletRepo {
	return &WalletRepo{pool: pool}
}

func (r *WalletRepo) IncreaseBalance(ctx context.Context, userId string, walletId string, amount float64, currency string) error {
	_, err := r.pool.Exec(ctx, increaseBalanceQuery, walletId, userId, amount, currency)
	return err
}

func (r *WalletRepo) DecreaseBalance(ctx context.Context, userId string, walletId string, amount float64, currency string) error {
	_, err := r.pool.Exec(ctx, decreaseBalanceQuery, walletId, userId, amount, currency)
	return err
}
