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

CREATE TABLE "hamonix_business"."users" (
  "id" varchar PRIMARY KEY,
  "name" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "hamonix_business"."u_portfolio" (
  "id" bigserial PRIMARY KEY,
  "user_id" varchar NOT NULL,
  "portfolio_id" varchar
);

CREATE TABLE "hamonix_business"."portfolios" (
  "id" varchar PRIMARY KEY,
  "name" varchar UNIQUE NOT NULL,
  "privacy" varchar NOT NULL,
  "author_id" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "hamonix_business"."assets" (
  "id" bigserial PRIMARY KEY,
  "portfolio_id" varchar NOT NULL,
  "ticker_id" int NOT NULL,
  "price" float NOT NULL,
  "allocation" float NOT NULL
);

CREATE TABLE "hamonix_business"."ticker_prices" (
  "ticker_id" bigserial PRIMARY KEY,
  "open" float NOT NULL,
  "high" float NOT NULL,
  "low" float NOT NULL,
  "close" float NOT NULL,
  "date" date NOT NULL
);

CREATE TABLE "hamonix_business"."portfolio_categories" (
  "id" varchar PRIMARY KEY,
  "name" varchar NOT NULL,
  "description" varchar,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "hamonix_business"."u_categories" (
  "id" bigserial PRIMARY KEY,
  "category_id" varchar,
  "user_id" varchar NOT NULL
);

CREATE TABLE "hamonix_business"."p_categories" (
  "id" bigserial PRIMARY KEY,
  "portfolio_id" varchar NOT NULL,
  "category_id" varchar
);

CREATE TABLE "hamonix_business"."p_branches" (
  "id" bigserial PRIMARY KEY,
  "portfolio_id" varchar NOT NULL,
  "branch_id" varchar
);

CREATE TABLE "hamonix_business"."p_advisors" (
  "id" bigserial PRIMARY KEY,
  "portfolio_id" varchar NOT NULL,
  "advisor_id" varchar
);

CREATE TABLE "hamonix_business"."p_organizations" (
  "id" bigserial PRIMARY KEY,
  "portfolio_id" varchar NOT NULL,
  "organization_id" varchar
);

CREATE TABLE "hamonix_business"."eq_whitelables" (
  "id" varchar PRIMARY KEY,
  "name" varchar NOT NULL,
  "url" varchar NOT NULL,
  "description" varchar
);

CREATE TABLE "hamonix_business"."eq_backoffices" (
  "id" varchar PRIMARY KEY,
  "whitelable_id" varchar,
  "name" varchar NOT NULL,
  "description" varchar
);

CREATE TABLE "hamonix_business"."eq_organizations" (
  "id" varchar PRIMARY KEY,
  "backoffice_id" varchar,
  "code" varchar NOT NULL,
  "description" varchar
);

CREATE TABLE "hamonix_business"."eq_branchs" (
  "id" varchar PRIMARY KEY,
  "code" varchar NOT NULL,
  "description" varchar
);

CREATE TABLE "hamonix_business"."eq_advisors" (
  "id" varchar PRIMARY KEY,
  "code" varchar,
  "description" varchar
);

CREATE TABLE "hamonix_business"."eq_accounts" (
  "id" varchar PRIMARY KEY,
  "advisor_id" varchar,
  "code" varchar NOT NULL
);

CREATE INDEX ON "hamonix_business"."portfolios" USING BTREE ("name");
CREATE INDEX "created_at_index" ON "hamonix_business"."portfolios" ("created_at");
ALTER TABLE hamonix_business.tickers ADD CONSTRAINT unique_symbol_exchange UNIQUE (symbol, exchange);
CREATE INDEX idx_symbol_exchange ON hamonix_business.tickers (symbol, exchange);
CREATE INDEX idx_ticker_id ON hamonix_business.tickers (id);
