CREATE TABLE IF NOT EXISTS orders (
    order_uid TEXT PRIMARY KEY,
    track_number TEXT,
    entry TEXT,
    locale TEXT,
    internal_signature TEXT,
    customer_id TEXT,
    delivery_service TEXT,
    shardkey TEXT,
    sm_id INTEGER,
    date_created TEXT,
    oof_shard TEXT
);

CREATE TABLE IF NOT EXISTS delivery (
    order_uid TEXT PRIMARY KEY,
    name TEXT,
    phone TEXT,
    zip TEXT,
    city TEXT,
    address TEXT,
    region TEXT,
    email TEXT,
    FOREIGN KEY (order_uid) REFERENCES orders(order_uid)
);

CREATE TABLE IF NOT EXISTS payment (
    transaction TEXT PRIMARY KEY,
    order_uid TEXT,
    request_id TEXT,
    currency TEXT,
    provider TEXT,
    amount INTEGER,
    payment_dt BIGINT,
    bank TEXT,
    delivery_cost INTEGER,
    goods_total INTEGER,
    custom_fee INTEGER,
    FOREIGN KEY (order_uid) REFERENCES orders(order_uid)
);

CREATE TABLE IF NOT EXISTS items (
    id SERIAL PRIMARY KEY,
    order_uid TEXT,
    chrt_id BIGINT,
    track_number TEXT,
    price INTEGER,
    rid TEXT,
    name TEXT,
    sale INTEGER,
    size TEXT,
    total_price INTEGER,
    nm_id INTEGER,
    brand TEXT,
    status INTEGER,
    FOREIGN KEY (order_uid) REFERENCES orders(order_uid)
);

INSERT INTO orders (
    order_uid, track_number, entry, locale, internal_signature, customer_id,
    delivery_service, shardkey, sm_id, date_created, oof_shard
) VALUES (
    'b563feb7b2b84b6test', 'WBILMTESTTRACK', 'WBIL', 'en', '',
    'test', 'meest', '9', 99, '2021-11-26T06:22:19Z', '1'
)
ON CONFLICT (order_uid) DO NOTHING;

INSERT INTO delivery (
    order_uid, name, phone, zip, city, address, region, email
) VALUES (
    'b563feb7b2b84b6test', 'Test Testov', '+9720000000', '2639809',
    'Kiryat Mozkin', 'Ploshad Mira 15', 'Kraiot', 'test@gmail.com'
)
ON CONFLICT (order_uid) DO NOTHING;

INSERT INTO payment (
    transaction, order_uid, request_id, currency, provider, amount,
    payment_dt, bank, delivery_cost, goods_total, custom_fee
) VALUES (
    'b563feb7b2b84b6test', 'b563feb7b2b84b6test', '', 'USD', 'wbpay',
    1817, 1637907727, 'alpha', 1500, 317, 0
)
ON CONFLICT (transaction) DO NOTHING;

INSERT INTO items (
    order_uid, chrt_id, track_number, price, rid, name,
    sale, size, total_price, nm_id, brand, status
) VALUES (
    'b563feb7b2b84b6test', 9934930, 'WBILMTESTTRACK', 453,
    'ab4219087a764ae0btest', 'Mascaras', 30, '0', 317, 2389212,
    'Vivienne Sabo', 202
)
ON CONFLICT (id) DO NOTHING;
