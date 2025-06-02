CREATE TABLE IF NOT EXISTS user_invitations (
    token bytea,
    user_id bigint NOT NULL,
    PRIMARY key (token, user_id)
)