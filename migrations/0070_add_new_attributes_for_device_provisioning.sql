-- +migrate Up

alter table device
    add column provision_id char(20) unique;
alter table device
    add column model text;
alter table device
    add column serial_number text;
alter table device
    add column manufacturer text;

alter table gateway_profile
    add constraint gateway_profile_name unique (name);

alter table service_profile
    add constraint service_profile_name unique (name);
