-- +migrate Up
alter table gateway alter COLUMN model set not null;
alter table gateway alter COLUMN model set default '';
update table gateway set model= '' where model = 'null';

-- +migrate Down
