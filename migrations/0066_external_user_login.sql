-- +migrate Up

create table external_login
(
    user_id     bigint not null references "user" on delete cascade,
    service     text   not null,
    external_id text   not null
);