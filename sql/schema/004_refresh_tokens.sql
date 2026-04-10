-- +goose Up

CREATE TABLE refresh_tokens (
    token text PRIMARY KEY,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    expires_at timestamptz NOT NULL,
    revoked_at timestamptz DEFAULT NULL,
    user_id uuid NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down

DROP TABLE refresh_tokens;
