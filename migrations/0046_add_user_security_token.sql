-- +migrate Up
alter table "user" add column security_token character varying(200) null;
alter table "user" alter column password_hash drop not null;
update "user" set password_hash = null where password_hash = '';

-- +migrate Down
alter table "user" drop column security_token;
update "user" set password_hash = '' where password_hash is null;
alter table "user" alter column password_hash set not null;
