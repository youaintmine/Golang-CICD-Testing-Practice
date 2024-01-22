CREATE TABLE "entries" (
  "id" bigserial PRIMARY KEY,
  "account_id" bigint NOT NULL,
  "amount" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "accounts" (
  "id" bigserial PRIMARY KEY,
  "owner" varchar NOT NULL,
  "balance" bigint NOT NULL,
  "currency" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT 'now()',
  "country_code" int
);

CREATE TABLE "transfers" (
  "id" bigserial PRIMARY KEY,
  "from_acc" bigint NOT NULL,
  "to_acc" bigint NOT NULL,
  "amount" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE INDEX ON "entries" ("account_id");

CREATE INDEX ON "accounts" ("owner");

CREATE INDEX ON "transfers" ("from_acc");

CREATE INDEX ON "transfers" ("to_acc");

CREATE INDEX ON "transfers" ("from_acc", "to_acc");

COMMENT ON COLUMN "entries"."amount" IS 'can be positive or negative';

COMMENT ON COLUMN "transfers"."amount" IS 'needs to be positive';

ALTER TABLE "entries" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("from_acc") REFERENCES "accounts" ("id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("to_acc") REFERENCES "accounts" ("id");
