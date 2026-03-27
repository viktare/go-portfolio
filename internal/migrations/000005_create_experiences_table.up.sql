CREATE TABLE IF NOT EXISTS experiences (
    id          VARCHAR(36)  PRIMARY KEY DEFAULT gen_random_uuid(),
    role        VARCHAR(100) NOT NULL,
    description TEXT         NOT NULL,
    company_id  VARCHAR(36)  NOT NULL,
    start_date  DATE         NOT NULL,
    end_date    DATE,
    order_index INT          NOT NULL DEFAULT 0,

    CONSTRAINT fk_experiences_company
        FOREIGN KEY (company_id)
        REFERENCES companies(id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE
);