START TRANSACTION;

-- add tier in organization
ALTER TABLE "public"."organizations" ADD COLUMN "tier" VARCHAR(255) NOT NULL DEFAULT 'free';

COMMIT;