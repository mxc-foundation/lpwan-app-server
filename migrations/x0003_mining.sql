-- +migrate Up
alter table gateway
    add column model varchar(32),
    add column first_heartbeat int not null default 0,
    add column last_heartbeat int not null default 0;

alter table gateway
    alter column first_heartbeat drop default,
    alter column last_heartbeat drop default;

-- +migrate Down
alter table gateway
    drop column model,
    drop column first_heartbeat,
    drop column last_heartbeat;
