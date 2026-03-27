CREATE TABLE IF NOT EXISTS technology_fields (
    id   VARCHAR(36)  PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    code VARCHAR(100) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS technologies (
    id       VARCHAR(36)  PRIMARY KEY DEFAULT gen_random_uuid(),
    name     VARCHAR(100) NOT NULL,
    code     VARCHAR(100) NOT NULL UNIQUE,
    field_id VARCHAR(36)  NOT NULL,

    CONSTRAINT fk_technologies_fields
        FOREIGN KEY (field_id)
        REFERENCES technology_fields(id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE
);