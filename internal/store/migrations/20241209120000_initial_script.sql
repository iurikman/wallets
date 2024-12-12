-- +migrate Up

CREATE TABLE wallets (
    id uuid primary key,
    balance numeric not null DEFAULT 0 check (balance >= 0),
    created_at timestamp not null,
    updated_at timestamp,
    deleted bool not null
);

CREATE TABLE transactions_history (
  id uuid primary key,
  wallet_id uuid references wallets (id),
  amount numeric not null,
  transaction_type varchar not null,
  executed_at timestamp not null
);

CREATE INDEX idx_balance ON wallets (balance);
CREATE INDEX idx_wallet_id ON transactions_history (wallet_id);
CREATE INDEX idx_executed_at ON transactions_history (executed_at);

-- +migrate Down

DROP TABLE wallets, transactions_history CASCADE;
