-- Cleanup
DROP TABLE IF EXISTS users, items, item_prices, user_items, sessions CASCADE;

-- uuid support
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Tables
CREATE TABLE IF NOT EXISTS users (
    id uuid NOT NULL DEFAULT uuid_generate_v4 (),
    username text NOT NULL,
    email text NOT NULL UNIQUE,
    password TEXT NOT NULL,
    created timestamptz NOT NULL DEFAULT now(),
    logged_in timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS items (
    id uuid NOT NULL DEFAULT uuid_generate_v4 (),
    name text NOT NULL,
    description text,
    image_url text NOT NULL,
    url text NOT NULL,
    currency text NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS item_prices (
    item_id uuid NOT NULL REFERENCES items (id) ON DELETE CASCADE,
    time timestamptz NOT NULL DEFAULT NOW(),
    price int,
    available boolean DEFAULT TRUE,
    PRIMARY KEY (item_id, time)
);

CREATE TABLE IF NOT EXISTS user_items (
    user_id uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    item_id uuid NOT NULL REFERENCES items (id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, item_id)
);

CREATE TABLE IF NOT EXISTS sessions (
    id uuid NOT NULL DEFAULT uuid_generate_v4 (),
    user_id uuid REFERENCES users (id) ON DELETE CASCADE,
    expires_after timestamptz NOT NULL,
    jwt text NOT NULL,
    PRIMARY KEY (id, user_id)
);

CREATE TABLE IF NOT EXISTS subscriptions (
    user_id uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    item_id uuid NOT NULL REFERENCES items (id) ON DELETE CASCADE,
    email text NOT NULL,
    target_price int,
    PRIMARY KEY (user_id, item_id)
);

