ALTER TABLE IF EXISTS accounts DROP CONSTRAINT IF EXISTS "fk_account_owner";

ALTER TABLE IF EXISTS accounts DROP CONSTRAINT IF EXISTS "unique_owner_currency";  

DROP TABLE IF EXISTS "users";