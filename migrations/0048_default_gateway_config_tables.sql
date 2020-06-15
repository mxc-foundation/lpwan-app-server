-- +migrate Up
CREATE TABLE default_gateway_config (
    id bigserial,
    model varchar(32) not null default '',
    region varchar(32) not null default '',
    created_at timestamp with time zone not null,
    updated_at timestamp with time zone not null,
    default_config text not null default '',
    PRIMARY KEY (id)
);

-- +migrate Down
