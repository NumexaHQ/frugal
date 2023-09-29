START TRANSACTION;

-- remove tier in organization
ALTER TABLE "public"."organizations" DROP COLUMN "tier";

COMMIT;