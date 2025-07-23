INSERT INTO wallet_balance (walletId, currency, amount) 
VALUES ($1, $2, $3)
ON CONFLICT (currency) DO UPDATE
SET amount = wallet_balance.amount + EXCLUDED.amount
    WHERE EXCLUDED.amount > 0
        AND wallet_balance.walletId = $1
        AND wallet_balance.currency = $2;
