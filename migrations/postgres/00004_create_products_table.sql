-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS products(
    id SERIAL PRIMARY KEY,
    sku VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(500) NOT NULL,
    description TEXT,
    category_id INTEGER REFERENCES categories(id) ON DELETE SET NULL,
    manufacturer_id INTEGER REFERENCES manufacturers(id) ON DELETE SET NULL,
    weight INTEGER NOT NULL,
    unit VARCHAR(20) DEFAULT 'pc.',
    price INTEGER NOT NULL,     
    baseprice INTEGER NOT NULL,
    issueyear SMALLINT NOT NULL CHECK(issueyear BETWEEN 1900 AND 2100),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS products;
-- +goose StatementEnd

