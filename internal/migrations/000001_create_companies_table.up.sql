create table if not exists companies (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    address VARCHAR(255) NOT NULL,
    website VARCHAR(255) DEFAULT NULL,
    logo VARCHAR(255) DEFAULT NULL
);