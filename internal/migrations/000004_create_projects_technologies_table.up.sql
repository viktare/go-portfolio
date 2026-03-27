CREATE TABLE IF NOT EXISTS project_technologies (
    project_id    VARCHAR(36) NOT NULL,
    technology_id VARCHAR(36) NOT NULL,

    PRIMARY KEY (project_id, technology_id),

    CONSTRAINT fk_project_technologies_project
        FOREIGN KEY (project_id)
        REFERENCES projects(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT fk_project_technologies_technology
        FOREIGN KEY (technology_id)
        REFERENCES technologies(id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE
);