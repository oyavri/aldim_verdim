package frontend

import (
	"context"
	_ "embed"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oyavri/aldim_verdim/pkg/entity"
)

var (
	//go:embed sql/get_balance.sql
	getBalanceQuery string
)

type WalletRepository interface {
	GetWallets(context.Context)
}

type WalletRepo struct {
	pool *pgxpool.Pool
}

func NewWalletRepository(pool *pgxpool.Pool) *WalletRepo {
	return &WalletRepo{pool: pool}
}

func (r *WalletRepo) GetWallets(ctx context.Context) ([]entity.Wallet, error) {
	rows, err := r.pool.Query(ctx, getBalanceQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	walletsMap := make(map[string]*entity.Wallet)

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
			walletsMap[walletId] = &entity.Wallet{
				Id:       walletId,
				Balances: []entity.Balance{},
			}
		}

		walletsMap[walletId].Balances = append(walletsMap[walletId].Balances, entity.Balance{
			Amount:   amount,
			Currency: currency,
		})
	}

	wallets := make([]entity.Wallet, 0, len(walletsMap))
	for _, wallet := range walletsMap {
		wallets = append(wallets, *wallet)
	}

	return wallets, nil
}
