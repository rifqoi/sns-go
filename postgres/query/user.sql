-- name: AddUser :one
INSERT INTO users (
    id,
    username,
    name,
    password,
    email,
    age,
    gender,
    date_of_birth,
    created_at,
    updated_at
)
VALUES (
    @id,
    @username,
    @name,
    @password,
    @email,
    @age,
    @gender,
    @date_of_birth,
    current_timestamp,
    current_timestamp
) RETURNING *;

-- name: GetUsers :many
SELECT id, 
    username,
    name,
    email,
    age,
    gender,
    date_of_birth,
    created_at,
    updated_at 
from users;

-- name: GetUserByID :one
SELECT id, 
    username,
    name,
    email,
    age,
    gender,
    date_of_birth,
    created_at,
    updated_at 
from users where id = @id;

-- name: GetUserByEmail :one
SELECT id, 
    username,
    name,
    email,
    age,
    gender,
    date_of_birth,
    created_at,
    updated_at 
from users where email = @email;
