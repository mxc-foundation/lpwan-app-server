-- +migrate Up

create table external_login
(
    user_id           bigint not null references "user" on delete cascade,
    service           text   not null,
    external_id       text   not null,
    external_username text   not null
);

alter table "user"
    drop column external_id;
alter table "user"
    drop column note;
alter table "user"
    drop column email_old;
alter table "user"
    add column display_name text;

update "user"
set display_name = "user".email;
