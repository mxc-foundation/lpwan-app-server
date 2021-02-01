-- +migrate Up

alter table external_login
    add column verification text not null default '';