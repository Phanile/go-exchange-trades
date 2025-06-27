-- +goose Up
-- +goose StatementBegin

CREATE TABLE Coins (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    shortname TEXT NOT NULL,
    chain_name TEXT,
    image_url TEXT
);

CREATE TABLE OrderTypes (
    id SMALLSERIAL PRIMARY KEY,
    title TEXT NOT NULL UNIQUE
);

CREATE TABLE OrderSides (
    id SMALLSERIAL PRIMARY KEY,
    title TEXT NOT NULL UNIQUE
);

CREATE TABLE Orders (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    send_coin_id BIGINT NOT NULL REFERENCES coins(id) ON DELETE RESTRICT,
    receive_coin_id BIGINT NOT NULL REFERENCES coins(id) ON DELETE RESTRICT,
    ordertype_id SMALLINT NOT NULL REFERENCES OrderTypes(id) ON DELETE RESTRICT,
    orderside_id SMALLINT NOT NULL REFERENCES OrderSides(id) ON DELETE RESTRICT,
    amount NUMERIC(30, 18) NOT NULL,
    price NUMERIC(30, 18),
    status TEXT NOT NULL DEFAULT 'open',
    timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_orders_user_id ON Orders(user_id);
CREATE INDEX idx_orders_pair ON Orders(send_coin_id, receive_coin_id);

INSERT INTO OrderTypes (id, title) VALUES
   (0, 'market'),
   (1, 'limit');

INSERT INTO OrderSides (id, title) VALUES
   (0, 'buy'),
   (1, 'sell');

CREATE TABLE Trades (
    id BIGSERIAL PRIMARY KEY,
    buy_order_id BIGINT NOT NULL REFERENCES Orders(id) ON DELETE CASCADE,
    sell_order_id BIGINT NOT NULL REFERENCES Orders(id) ON DELETE CASCADE,
    amount NUMERIC(30, 18) NOT NULL,
    price NUMERIC(30, 18) NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_trades_buy_order_id ON Trades(buy_order_id);
CREATE INDEX idx_trades_sell_order_id ON Trades(sell_order_id);
CREATE INDEX idx_trades_pair_price ON Trades(price);

INSERT INTO Coins (name, shortname, chain_name, image_url) VALUES
   ('Solana', 'SOL', 'Solana', 'https://example.com/solana.png'),
   ('TON', 'TON', 'TON', 'https://example.com/ton.png'),
   ('SCMN', 'SCMN', 'uretra_network', 'https://example.com/scmn.png');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS Orders;
DROP TABLE IF EXISTS OrderSides;
DROP TABLE IF EXISTS OrderTypes;
DROP TABLE IF EXISTS Coins;
DROP TABLE IF EXISTS Trades;

-- +goose StatementEnd
