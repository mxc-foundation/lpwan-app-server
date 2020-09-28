-- +migrate Up
CREATE TABLE password_reset (
    user_id bigint REFERENCES "user" (id) ON DELETE CASCADE,
    otp text NOT NULL,
    generated_at timestamp NOT NULL,
    attempts_left int NOT NULL,
    PRIMARY KEY (user_id)
);

-- +migrate Down
DROP TABLE password_reset;
