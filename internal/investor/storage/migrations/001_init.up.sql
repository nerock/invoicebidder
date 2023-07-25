CREATE TYPE amount AS (
   number NUMERIC,
   currency_code TEXT
);

CREATE TABLE investors (
   id CHAR(36) PRIMARY KEY,
   name TEXT NOT NULL
);

CREATE TABLE investor_balances (
    investor_id CHAR(36) NOT NULL REFERENCES investors (id),
    amount amount NOT NULL
);