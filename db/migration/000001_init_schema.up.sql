CREATE SCHEMA IF NOT EXISTS hamonix_business;
CREATE TABLE hamonix_business.tickers (
  "id" bigserial PRIMARY KEY,
  "symbol" varchar NOT NULL,
  "description" varchar NOT NULL,
  "exchange" varchar NOT NULL,
  "currency" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE hamonix_business.tickers ADD CONSTRAINT unique_symbol_exchange UNIQUE (symbol, exchange);
CREATE INDEX idx_symbol_exchange ON hamonix_business.tickers (symbol, exchange);
CREATE INDEX idx_ticker_id ON hamonix_business.tickers (id);
