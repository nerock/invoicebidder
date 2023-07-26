DROP TABLE issuer_balances;

ALTER TABLE issuers
ADD COLUMN balance amount NOT NULL