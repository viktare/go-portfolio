CREATE TABLE IF NOT EXISTS experience_technologies (
    experience_id VARCHAR(36) NOT NULL,
    technology_id VARCHAR(36) NOT NULL,

    PRIMARY KEY (experience_id, technology_id),

    CONSTRAINT fk_experience_technologies_experience
        FOREIGN KEY (experience_id)
        REFERENCES experiences(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT fk_experience_technologies_technology
        FOREIGN KEY (technology_id)
        REFERENCES technologies(id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE
);