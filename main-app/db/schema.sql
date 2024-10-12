CREATE TYPE role_enum AS ENUM ('inviter', 'invitee');

CREATE TABLE connections (
    connection_id VARCHAR NOT NULL,
    id BIGSERIAL NOT NULL,
    alias VARCHAR NOT NULL,
    my_role role_enum NOT NULL,
    PRIMARY KEY (connection_id),
    FOREIGN KEY (id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS schemas (
    schema_id VARCHAR NOT NULL,
    credential_definition_id VARCHAR NOT NULL,
    schema_name VARCHAR NOT NULL,
    PRIMARY KEY (schema_id)
);
