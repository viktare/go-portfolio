CREATE TABLE IF NOT EXISTS files (
    id        VARCHAR(36)  PRIMARY KEY DEFAULT gen_random_uuid(),
    name      VARCHAR(255) NOT NULL,
    path      VARCHAR(500) NOT NULL,
    mime_type VARCHAR(100) NOT NULL,
    size      BIGINT       NOT NULL
);