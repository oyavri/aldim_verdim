INSERT INTO wallet_balance (walletId, currency, amount) 
VALUES ($1, $2, $3)
ON CONFLICT (walletId, currency) DO UPDATE
SET amount = wallet_balance.amount + EXCLUDED.amount
    WHERE EXCLUDED.amount > 0
