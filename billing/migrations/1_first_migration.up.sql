CREATE
EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE currencies
(
    id         UUID PRIMARY KEY     DEFAULT uuid_generate_v4(),
    code       TEXT        NOT NULL UNIQUE,
    name       TEXT        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

INSERT INTO currencies (code, name)
VALUES ('USD', 'United States Dollar'),
       ('EUR', 'Euro'),
       ('GBP', 'British Pound Sterling'),
       ('JPY', 'Japanese Yen'),
       ('AUD', 'Australian Dollar'),
       ('CAD', 'Canadian Dollar'),
       ('CHF', 'Swiss Franc'),
       ('CNY', 'Chinese Yuan'),
       ('SEK', 'Swedish Krona'),
       (code 'GEL', name 'Georgian Lari'),
       ('NZD', 'New Zealand Dollar');

CREATE TYPE bill_status AS ENUM ('OPEN', 'CLOSED');

CREATE TABLE bills
(
    id           UUID PRIMARY KEY        DEFAULT uuid_generate_v4(),
    customer_id  TEXT           NOT NULL,
    currency_id  UUID           NOT NULL REFERENCES currencies (id),
    amount       NUMERIC(20, 6) NOT NULL,
    status       bill_status    NOT NULL DEFAULT 'OPEN',
    period_start TIMESTAMPTZ    NOT NULL,
    period_end   TIMESTAMPTZ,
    created_at   TIMESTAMPTZ    NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ    NOT NULL DEFAULT now(),
    deleted_at   TIMESTAMPTZ,
    closed_at    TIMESTAMPTZ
);

CREATE TABLE bill_items
(
    id          UUID PRIMARY KEY        DEFAULT uuid_generate_v4(),
    bill_id     UUID           NOT NULL REFERENCES bills (id) ON DELETE CASCADE,
    description TEXT           NOT NULL,
    amount      NUMERIC(20, 6) NOT NULL,
    created_at  TIMESTAMPTZ    NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ    NOT NULL DEFAULT now(),
    deleted_at  TIMESTAMPTZ
);