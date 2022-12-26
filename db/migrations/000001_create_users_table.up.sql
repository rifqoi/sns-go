-- Create enum for gender
CREATE TYPE gender AS ENUM ('Male', 'Female');

CREATE TABLE IF NOT EXISTS users(
    id uuid PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(300) NOT NULL,
    name varchar(50) NOT NULL,
    gender gender not null,
    email VARCHAR(50) UNIQUE NOT NULL,
    age INT NOT NULL,
    date_of_birth TIMESTAMPTZ not null,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create function to update timestamp on field updated_at
CREATE FUNCTION update_updated_at_users()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create trigger to trigger the function
CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE
    ON
        users
    FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_users();
