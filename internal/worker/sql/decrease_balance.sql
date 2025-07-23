UPDATE wallet_balance 
SET amount = wallet_balance.amount - $3
WHERE $3 > 0
    AND walletId = $1
    AND currency = $2;
