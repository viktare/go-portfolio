create table if not exists companies (
    id      VARCHAR(36)  PRIMARY KEY DEFAULT gen_random_uuid(),
    name    VARCHAR(255) NOT NULL,
    address VARCHAR(255) NOT NULL,
    website VARCHAR(255) DEFAULT NULL,
    logo    VARCHAR(255) DEFAULT NULL
);