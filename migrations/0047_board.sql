-- +migrate Up
create table board
(
    mac bytea NOT NULL,
    sn character varying(16) COLLATE pg_catalog."default",
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    model character varying(16) COLLATE pg_catalog."default" NOT NULL,
    vpn_addr inet NOT NULL,
    qa_err integer NOT NULL,
    os_version character varying(16) COLLATE pg_catalog."default",
    fpga_version character varying(16) COLLATE pg_catalog."default",
    root_password character varying(64),
    server character varying(64),
    CONSTRAINT board_pkey PRIMARY KEY (mac),
    CONSTRAINT board_sn_key UNIQUE (sn)
);

-- +migrate Down
drop table board;
