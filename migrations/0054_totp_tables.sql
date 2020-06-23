-- +migrate Up
CREATE TABLE totp_configuration (
    user_id bigint REFERENCES "user" (id) ON DELETE CASCADE,
    is_enabled boolean NOT NULL DEFAULT false,
    secret text NOT NULL,
    last_time_slot bigint NOT NULL DEFAULT 0,
    PRIMARY KEY (user_id)
);
CREATE TABLE totp_recovery_codes (
    id SERIAL,
    user_id bigint NOT NULL REFERENCES "user" (id) ON DELETE CASCADE,
    code text NOT NULL,
    PRIMARY KEY (id)
);


-- +migrate Down
DROP TABLE totp_configuration;
DROP TABLE totp_recovery_codes;
