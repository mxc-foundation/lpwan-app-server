-- +migrate Up

alter table device
    add column provision_id char(20) unique;
alter table device
    add column model text;
alter table device
    add column serial_number text;
alter table device
    add column manufacturer text;