CREATE TYPE user_role AS ENUM ('customer','admin')

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(100) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    phone VARCHAR(20),
    is_active BOOLEAN DEFAULT true,
    role user_role DEFAULT 'customer',
    created_at TIME WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIME WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIME WITH TIME ZONE 

);

CREATE INDEX idx_users_email on users(email);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);