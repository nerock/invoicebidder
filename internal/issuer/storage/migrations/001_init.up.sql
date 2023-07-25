CREATE TYPE amount AS (
   number NUMERIC,
   currency_code TEXT
);

CREATE TABLE issuers (
   id CHAR(36) PRIMARY KEY,
   name TEXT NOT NULL
);

CREATE TABLE issuer_balances (
    issuer_id CHAR(36) NOT NULL REFERENCES issuers (id),
    amount amount NOT NULL
)