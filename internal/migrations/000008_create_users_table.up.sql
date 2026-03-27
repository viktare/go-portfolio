CREATE TABLE IF NOT EXISTS users (
    id             VARCHAR(36)  PRIMARY KEY DEFAULT gen_random_uuid(),
    name           VARCHAR(100) NOT NULL,
    surname        VARCHAR(100) NOT NULL,
    bio            TEXT         NOT NULL,
    birth_date     DATE,
    location       VARCHAR(255) NOT NULL,
    resume_file_id VARCHAR(36),

    CONSTRAINT fk_users_resume_file
        FOREIGN KEY (resume_file_id)
        REFERENCES files(id)
        ON DELETE SET NULL
        ON UPDATE CASCADE
);