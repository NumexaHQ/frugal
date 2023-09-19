START TRANSACTION;

DROP TABLE IF EXISTS "public"."provider_keys";

DROP TABLE IF EXISTS "public"."provider_secrets";

DROP TABLE IF EXISTS "public"."nxa_api_key_property";

ALTER TABLE "public"."nxa_api_key" DROP COLUMN "nxa_api_key_property_id";

COMMIT;