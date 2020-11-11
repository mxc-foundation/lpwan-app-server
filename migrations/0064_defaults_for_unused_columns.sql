-- +migrate Up
ALTER TABLE "user" ALTER COLUMN session_ttl SET DEFAULT 0;
ALTER TABLE "user" ALTER COLUMN "note" SET DEFAULT '';
