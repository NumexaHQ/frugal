START TRANSACTION;

CREATE TABLE IF NOT EXISTS "public"."setting" (
    "id" SERIAL PRIMARY KEY,
    "key" VARCHAR(255) NOT NULL UNIQUE,
    "value" JSONB NOT NULL,
    "visible" BOOLEAN NOT NULL DEFAULT FALSE,
    "created_at" TIMESTAMP NOT NULL,
    "updated_at" TIMESTAMP NOT NULL
);


CREATE TABLE IF NOT EXISTS "public"."provider_keys" (
    "id" SERIAL PRIMARY KEY,
    "key_uuid" VARCHAR(255) NOT NULL UNIQUE,
    "name"  VARCHAR(255) NOT NULL, 
    "provider" VARCHAR(255) NOT NULL,
    "creator_id" INTEGER NOT NULL REFERENCES users(id),
    "organization_id" INTEGER NOT NULL REFERENCES organizations(id),
    "project_id" INTEGER NOT NULL REFERENCES projects(id),
    "created_at" TIMESTAMP NOT NULL,
    "updated_at" TIMESTAMP NOT NULL 
);

CREATE TABLE IF NOT EXISTS "public"."provider_secrets" (
    "id" SERIAL PRIMARY KEY,
    "provider_key_id" INTEGER NOT NULL REFERENCES provider_keys(id),
    "type" VARCHAR(255) NOT NULL,
    "key" VARCHAR(255) NOT NULL,
    "created_at" TIMESTAMP NOT NULL,
    "updated_at" TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS "public"."nxa_api_key_property" (
    "id" SERIAL PRIMARY KEY,
    "rate_limit" INTEGER NOT NULL,
    "rate_limit_period" VARCHAR(255) NOT NULL,
    "enforce_caching" BOOLEAN NOT NULL,
    "overall_cost_limit" INTEGER NOT NULL,
    "alert_on_threshold" INTEGER NOT NULL,
    "provider_key_id" INTEGER NULL REFERENCES provider_keys(id),
    "expires_at" TIMESTAMP NOT NULL,  
    "created_at" TIMESTAMP NOT NULL,
    "updated_at" TIMESTAMP NOT NULL
);

ALTER TABLE "public"."nxa_api_key" ADD COLUMN "nxa_api_key_property_id" INTEGER NULL REFERENCES nxa_api_key_property(id);
ALTER TABLE "public"."nxa_api_key" ADD COLUMN "provider_key_id" INTEGER NULL REFERENCES provider_keys(id);
ALTER TABLE "public"."nxa_api_key" ADD COLUMN "revoked" BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE "public"."nxa_api_key" ADD COLUMN "revoked_at" TIMESTAMP NULL;
ALTER TABLE "public"."nxa_api_key" ADD COLUMN "revoked_by" INTEGER NULL REFERENCES users(id);
ALTER TABLE "public"."nxa_api_key" ADD COLUMN "disabled" BOOLEAN NOT NULL DEFAULT FALSE;

COMMIT;