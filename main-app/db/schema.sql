CREATE TYPE role_enum AS ENUM ('inviter', 'invitee');
CREATE TABLE connections (
    connection_id VARCHAR NOT NULL,
    id BIGSERIAL NOT NULL,
    alias VARCHAR,
    my_role role_enum NOT NULL,
    PRIMARY KEY (connection_id),
    FOREIGN KEY (id) REFERENCES users(id)
);

