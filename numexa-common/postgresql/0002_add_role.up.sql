START TRANSACTION;

INSERT INTO roles (name, created_at, updated_at) VALUES ('admin', now(), now()) ON CONFLICT (name) DO NOTHING;
INSERT INTO roles (name, created_at, updated_at) VALUES ('standard-user', now(), now()) ON CONFLICT (name) DO NOTHING;
INSERT INTO roles (name, created_at, updated_at) VALUES ('read-only-user', now(), now()) ON CONFLICT (name) DO NOTHING;

COMMIT;