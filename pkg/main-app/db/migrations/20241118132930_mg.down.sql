-- Drop the connections table first because it references the users table
DROP TABLE IF EXISTS connections;

-- Drop the users table
DROP TABLE IF EXISTS users;

-- Drop the schemas table
DROP TABLE IF EXISTS schemas;

-- Drop the custom type role_enum
DROP TYPE IF EXISTS role_enum;
