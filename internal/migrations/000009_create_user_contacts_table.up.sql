CREATE TABLE IF NOT EXISTS user_contacts (
    id      VARCHAR(36)  PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(36)  NOT NULL,
    name    VARCHAR(100) NOT NULL,
    code    VARCHAR(100) NOT NULL,
    link    VARCHAR(500) NOT NULL,

    CONSTRAINT fk_user_contacts_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);