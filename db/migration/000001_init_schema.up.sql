CREATE SCHEMA IF NOT EXISTS harmonix_business;

CREATE TABLE "harmonix_business"."portfolio_profiles" (
  "id" varchar PRIMARY KEY,
  "name" varchar NOT NULL,
  "privacy" varchar NOT NULL,
  "author_id" varchar NOT NULL,
  "advisors" text[],
  "branches" text[],
  "organizations" text[],
  "accounts" text[],
  "expected_return" float NOT NULL,
  "is_new_buy_point" BOOLEAN NOT NULL DEFAULT FALSE,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "harmonix_business"."assets" (
  "id" bigserial PRIMARY KEY,
  "portfolio_profile_id" varchar NOT NULL,
  "ticker_name" varchar NOT NULL,
  "price" float NOT NULL,
  "allocation" float NOT NULL
);

CREATE TABLE "harmonix_business"."hrn_profile_account" (
  "id" bigserial PRIMARY KEY,
  "profile_id" varchar NOT NULL,
  "account_id" varchar NOT NULL,
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "harmonix_business"."portfolio_profiles" USING BTREE ("name");

CREATE INDEX "created_at_index" ON "harmonix_business"."portfolio_profiles" ("created_at");

-- ALTER TABLE "harmonix_business"."assets" ADD FOREIGN KEY ("portfolio_profile_id") REFERENCES "harmonix_business"."portfolio_profiles" ("id");
