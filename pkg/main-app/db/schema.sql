CREATE TYPE role_enum AS ENUM ('inviter', 'invitee');

CREATE TABLE IF NOT EXISTS connections (
    connection_id VARCHAR NOT NULL,
    id BIGSERIAL NOT NULL,
    my_mail_id TEXT,
    their_mail_id TEXT,
    PRIMARY KEY (connection_id),
    FOREIGN KEY (id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS schemas (
    schema_id VARCHAR NOT NULL,
    credential_definition_id VARCHAR NOT NULL,
    schema_name VARCHAR NOT NULL,
    attributes TEXT[],
    PRIMARY KEY (schema_id)
);
