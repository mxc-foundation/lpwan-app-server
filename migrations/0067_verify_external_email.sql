-- +migrate Up

alter table external_login
    add column verification text not null default '';

create table shopify_orders
(
    user_id      bigint                   not null references "user" on delete cascade,
    order_id     text                     not null,
    created_at   timestamp with time zone not null,
    product_id   text                     not null,
    order_amount integer                  not null,
    bonus_id     bigint                   not null default 0
);