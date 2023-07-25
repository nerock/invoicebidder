CREATE TYPE price AS (
   number NUMERIC,
   currency_code TEXT
);

CREATE TABLE invoices (
   id CHAR(36) PRIMARY KEY,
   issuer_id CHAR(36) NOT NULL,
   price price NOT NULL,
   status TEXT NOT NULL
);

CREATE TABLE bids (
    id CHAR(36) PRIMARY KEY,
    invoice_id CHAR(36) NOT NULL REFERENCES invoices (id),
    investor_id CHAR(36) NOT NULL,
    amount price NOT NULL,
    active BOOLEAN NOT NULL
);