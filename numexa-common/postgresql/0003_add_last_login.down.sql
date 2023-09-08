START TRANSACTION;

-- down of alter table users add column last_login timestamp;

ALTER TABLE "public"."users" DROP COLUMN last_login;
ALTER TABLE "public"."users" DROP COLUMN total_logins; 

COMMIT;