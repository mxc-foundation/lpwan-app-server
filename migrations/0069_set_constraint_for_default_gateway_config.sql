-- +migrate Up

alter table default_gateway_config
    add constraint model_region unique (model, region);