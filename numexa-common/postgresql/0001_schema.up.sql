START TRANSACTION;

CREATE TABLE IF NOT EXISTS "public"."organizations" (
    "id" SERIAL PRIMARY KEY,
    "name" VARCHAR(255) NOT NULL,
    "created_at" TIMESTAMP NOT NULL,
    "updated_at" TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS "public"."users" (
    "id" SERIAL PRIMARY KEY,
    "organization_id" INTEGER NOT NULL REFERENCES organizations(id),
    "name" VARCHAR(255) NOT NULL,
    "email" VARCHAR(255) NOT NULL UNIQUE,
    "password" VARCHAR(255) NOT NULL,
    "created_at" TIMESTAMP NOT NULL,
    "updated_at" TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS "public"."projects" (
    "id" SERIAL PRIMARY KEY,
    "organization_id" INTEGER NOT NULL REFERENCES organizations(id),
    "name" VARCHAR(255) NOT NULL,
    "description" TEXT
);

CREATE TABLE IF NOT EXISTS "public"."roles" (
    "id" SERIAL PRIMARY KEY,
    "name" VARCHAR(255) NOT NULL UNIQUE,
    "created_at" TIMESTAMP NOT NULL,
    "updated_at" TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS "public"."project_users" (
    "id" SERIAL PRIMARY KEY,
    "project_id" INTEGER NOT NULL REFERENCES projects(id),
    "user_id" INTEGER NOT NULL REFERENCES users(id),
    "role_id" INTEGER NOT NULL REFERENCES roles(id)
);

CREATE TABLE IF NOT EXISTS "public"."nxa_api_key" (
    "id" SERIAL PRIMARY KEY,
    "name" VARCHAR(255) NOT NULL,
    "api_key" VARCHAR(255) UNIQUE NOT NULL,
    "user_id" INTEGER NOT NULL REFERENCES users(id),
    "project_id" INTEGER NOT NULL REFERENCES projects(id),
    "created_at" TIMESTAMP NOT NULL,
    "updated_at" TIMESTAMP NOT NULL,
    "expires_at" TIMESTAMP NOT NULL
);

COMMIT;