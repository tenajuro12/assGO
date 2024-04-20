CREATE TABLE user_info (
                           id BIGSERIAL PRIMARY KEY,
                           created_at TIMESTAMP (0) WITH TIME ZONE NOT NULL DEFAULT now(),
                           updated_at TIMESTAMP (0) WITH TIME ZONE NOT NULL DEFAULT now(),
                           fname VARCHAR(255),
                           sname VARCHAR(255),
                           email VARCHAR(255)  UNIQUE NOT NULL,
                           password_hash BYTEA NOT NULL,
                           user_role VARCHAR(50),
                           activated BOOLEAN NOT NULL,
                           version INTEGER NOT NULL DEFAULT 1
);
