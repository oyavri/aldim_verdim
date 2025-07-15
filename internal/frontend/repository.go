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
	GetWallets(context.Context)
}

type WalletRepository struct {
	pool *pgxpool.Pool
}

func NewWalletRepository(pool *pgxpool.Pool) *WalletRepository {
	return &WalletRepository{pool: pool}
}

func (r *WalletRepository) GetWallets(ctx context.Context) ([]Wallet, error) {
	rows, err := r.pool.Query(ctx, getBalanceQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	walletsMap := make(map[string]*Wallet)

	for rows.Next() {
		var walletId string
		var amount float64
		var currency string

		err := rows.Scan(
			&walletId,
			&amount,
			&currency,
		)

		if err != nil {
			return nil, err
		}

		if _, exists := walletsMap[walletId]; !exists {
			walletsMap[walletId] = &Wallet{
				Id:       walletId,
				Balances: []Balance{},
			}
		}

		walletsMap[walletId].Balances = append(walletsMap[walletId].Balances, Balance{
			Amount:   amount,
			Currency: currency,
		})
	}

	wallets := make([]Wallet, 0, len(walletsMap))
	for _, wallet := range walletsMap {
		wallets = append(wallets, *wallet)
	}

	return wallets, nil
}
