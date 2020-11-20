-- +migrate Up
create table gateway_stc
(
    manufacturer_nr varchar(24) primary key,
    stc_org_id      bigint not null references organization on delete cascade
);

alter table gateway
    add column stc_org_id bigint references organization on delete cascade;
