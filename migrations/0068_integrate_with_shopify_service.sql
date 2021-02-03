-- +migrate Up

create table shopify_orders
(
    id                  bigserial primary key,
    org_id              bigint                   not null references "organization" on delete cascade,
    shopify_account_id  text                     not null,
    order_id            bigint                   not null,
    created_at          timestamp with time zone not null,
    product_id          bigint                   not null,
    amount_product      integer                  not null,
    bonus_per_piece_usd integer                  not null default 0,
    bonus_id            bigint                   not null default 0
);
