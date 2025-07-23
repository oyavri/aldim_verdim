SELECT w.id, b.currency, b.amount FROM wallet w
JOIN wallet_balance b 
    ON w.id = b.walletId;
