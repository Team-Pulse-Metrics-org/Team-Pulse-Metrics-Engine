-- Get user by ID
SELECT
    id,
    email,
    first_name,
    last_name,
    role,
    auth_provider,
    provider_data,
    created_at,
    updated_at
FROM users
WHERE id = $1;

-- Get user by email
SELECT
    id,
    email,
    first_name,
    last_name,
    role,
    auth_provider,
    provider_data,
    created_at,
    updated_at
FROM users
WHERE email = $1;

-- Get all users
SELECT
    id,
    email,
    first_name,
    last_name,
    role,
    auth_provider,
    provider_data,
    created_at,
    updated_at
FROM users;