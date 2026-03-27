CREATE TABLE IF NOT EXISTS technology_fields (
    id   SERIAL       PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    code VARCHAR(100) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS technologies (
    id       SERIAL       PRIMARY KEY,
    name     VARCHAR(100) NOT NULL,
    code     VARCHAR(100) NOT NULL UNIQUE,
    field_id INT          NOT NULL,

    CONSTRAINT fk_technologies_field
        FOREIGN KEY (field_id)
        REFERENCES technology_fields(id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE
);