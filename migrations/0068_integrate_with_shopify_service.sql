-- +migrate Up

create table shopify_orders
(
    user_id      bigint                   not null references "user" on delete cascade,
    order_id     text                     not null,
    created_at   timestamp with time zone not null,
    product_id   text                     not null,
    order_amount integer                  not null,
    bonus_id     bigint                   not null default 0
);
