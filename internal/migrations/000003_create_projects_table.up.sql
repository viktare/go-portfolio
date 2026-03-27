CREATE TABLE IF NOT EXISTS projects (
    id          VARCHAR(36)  PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(100) NOT NULL,
    code        VARCHAR(100) NOT NULL UNIQUE,
    status      VARCHAR(20)  NOT NULL CHECK (status IN ('in_progress', 'completed')),
    description TEXT,
    link        VARCHAR(255)
);
