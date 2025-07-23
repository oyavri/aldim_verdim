INSERT INTO wallet_balance (walletId, currency, amount) 
VALUES ($1, $2, $3)
ON CONFLICT (currency) DO UPDATE
SET amount = amount + EXCLUDED.amount
    WHERE EXCLUDED.amount > 0;
