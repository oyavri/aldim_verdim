-- To do

SELECT w.id FROM wallets w
JOIN balances b 
    ON w.id = b.id;
