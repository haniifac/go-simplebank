-- Create the accounts table
CREATE TABLE "accounts" (
  "id" BIGSERIAL PRIMARY KEY,
  "owner" VARCHAR NOT NULL,
  "balance" BIGINT NOT NULL,
  "currency" VARCHAR NOT NULL,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Create the entries table
CREATE TABLE "entries" (
  "id" BIGSERIAL PRIMARY KEY,
  "account_id" BIGINT NOT NULL,
  "amount" BIGINT NOT NULL,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Create the transfers table
CREATE TABLE "transfers" (
  "id" BIGSERIAL PRIMARY KEY,
  "from_account_id" BIGINT NOT NULL,
  "to_account_id" BIGINT NOT NULL,
  "amount" BIGINT NOT NULL,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Create indexes
CREATE INDEX "accounts_owner_idx" ON "accounts" ("owner");
CREATE INDEX "entries_account_id_idx" ON "entries" ("account_id");
CREATE INDEX "transfers_from_account_idx" ON "transfers" ("from_account_id");
CREATE INDEX "transfers_to_account_idx" ON "transfers" ("to_account_id");
CREATE INDEX "transfers_from_to_idx" ON "transfers" ("from_account_id", "to_account_id");

-- Add comments for clarity
COMMENT ON COLUMN "entries"."amount" IS 'can be negative or positive';
COMMENT ON COLUMN "transfers"."amount" IS 'must be positive';

-- Add foreign key constraints
ALTER TABLE "entries" ADD CONSTRAINT "entries_account_fk" FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");
ALTER TABLE "transfers" ADD CONSTRAINT "transfers_from_account_fk" FOREIGN KEY ("from_account_id") REFERENCES "accounts" ("id");
ALTER TABLE "transfers" ADD CONSTRAINT "transfers_to_account_fk" FOREIGN KEY ("to_account_id") REFERENCES "accounts" ("id");
ALTER TABLE "accounts" ADD CONSTRAINT "accounts_check_balance" CHECK ("balance" >= 0);
