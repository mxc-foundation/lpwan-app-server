-- +migrate Up

create table external_login
(
    user_id           bigint not null references "user" on delete cascade,
    service           text   not null,
    external_id       text   not null,
    external_username text   not null,
    constraint user_id_external_id unique (user_id, external_id)
);

alter table "user"
    drop column external_id;
alter table "user"
    drop column note;
alter table "user"
    drop column email_old;
alter table "user"
    add column display_name text not null default '';
alter table "user"
    add column last_login_service text not null default '';


update "user"
set display_name = email;
