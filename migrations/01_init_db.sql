DROP TABLE IF EXISTS wallet;
DROP TABLE IF EXISTS wallet_user;
DROP TABLE IF EXISTS user_wallets;
DROP TABLE IF EXISTS wallet_balance;

CREATE TABLE wallet (
    id VARCHAR PRIMARY KEY
);

CREATE TABLE wallet_user (
    id VARCHAR PRIMARY KEY
);

CREATE TABLE user_wallets (
    userId VARCHAR NOT NULL,
    walletId VARCHAR NOT NULL,
    FOREIGN KEY (userId) REFERENCES wallet_user(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    FOREIGN KEY (walletId) REFERENCES wallet(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);

CREATE TABLE wallet_balance (
    walletId VARCHAR NOT NULL,
    currency CHAR(3) UNIQUE,
    amount NUMERIC NOT NULL CHECK (amount >= 0),
    FOREIGN KEY (walletId) REFERENCES wallet(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);
