-- +migrate Up
alter table gateway
    add column config text not null default '',
    add column os_version varchar(16) not null default '',
    add column sn varchar(16) not null default '',
    add column statistics text not null default '',
    add column firmware_hash bytea not null default E'\\000';

create table gateway_firmware
(
    model         varchar(32) primary key ,
    resource_link varchar(256) not null default '',
    md5_hash      bytea not null default E'\\000'
);

insert into gateway_firmware
    (model)
    values ('MX1901'), ('MX1902'), ('MX1903');

-- +migrate Down

