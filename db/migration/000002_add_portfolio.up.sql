CREATE SCHEMA IF NOT EXISTS hamonix_business;

CREATE TYPE "portfolio_privacy" AS ENUM (
  'public',
  'private',
  'protected'
);

CREATE TABLE "hamonix_business"."portfolios" (
  "id" varchar PRIMARY KEY,
  "name" varchar UNIQUE NOT NULL,
  "privacy" portfolio_privacy NOT NULL,
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
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL,
  "description" varchar,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "hamonix_business"."p_categories" (
  "portfolio_id" varchar NOT NULL,
  "category_id" varchar
);

CREATE TABLE "hamonix_business"."p_branches" (
  "portfolio_id" varchar NOT NULL,
  "branch_id" varchar
);

CREATE TABLE "hamonix_business"."p_advisors" (
  "portfolio_id" varchar NOT NULL,
  "advisor_id" varchar
);

CREATE TABLE "hamonix_business"."p_organizations" (
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

-- CREATE INDEX "created_at_index" ON "hamonix_business"."portfolios" ("created_at");

-- ALTER TABLE "hamonix_business"."assets" ADD FOREIGN KEY ("portfolio_id") REFERENCES "hamonix_business"."portfolios" ("id");

-- ALTER TABLE "hamonix_business"."assets" ADD FOREIGN KEY ("ticker_id") REFERENCES "hamonix_business"."ticker_prices" ("ticker_id");

-- ALTER TABLE "hamonix_business"."p_categories" ADD FOREIGN KEY ("portfolio_id") REFERENCES "hamonix_business"."portfolios" ("id");

-- ALTER TABLE "hamonix_business"."p_categories" ADD FOREIGN KEY ("category_id") REFERENCES "hamonix_business"."portfolio_categories" ("id");

-- ALTER TABLE "hamonix_business"."p_branches" ADD FOREIGN KEY ("portfolio_id") REFERENCES "hamonix_business"."portfolios" ("id");

-- ALTER TABLE "hamonix_business"."p_branches" ADD FOREIGN KEY ("branch_id") REFERENCES "hamonix_business"."eq_branchs" ("id");

-- ALTER TABLE "hamonix_business"."p_advisors" ADD FOREIGN KEY ("portfolio_id") REFERENCES "hamonix_business"."portfolios" ("id");

-- ALTER TABLE "hamonix_business"."p_advisors" ADD FOREIGN KEY ("advisor_id") REFERENCES "hamonix_business"."eq_advisors" ("id");

-- ALTER TABLE "hamonix_business"."p_organizations" ADD FOREIGN KEY ("portfolio_id") REFERENCES "hamonix_business"."portfolios" ("id");

-- ALTER TABLE "hamonix_business"."p_organizations" ADD FOREIGN KEY ("organization_id") REFERENCES "hamonix_business"."eq_organizations" ("id");

-- ALTER TABLE "hamonix_business"."eq_backoffices" ADD FOREIGN KEY ("whitelable_id") REFERENCES "hamonix_business"."eq_whitelables" ("id");

-- ALTER TABLE "hamonix_business"."eq_organizations" ADD FOREIGN KEY ("backoffice_id") REFERENCES "hamonix_business"."eq_backoffices" ("id");

-- ALTER TABLE "hamonix_business"."eq_accounts" ADD FOREIGN KEY ("advisor_id") REFERENCES "hamonix_business"."eq_advisors" ("id");
