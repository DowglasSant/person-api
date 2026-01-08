-- Create operators table for authentication
CREATE TABLE IF NOT EXISTS operators (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,  -- Bcrypt hashed password
    active BOOLEAN DEFAULT true NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_operators_username ON operators(username);
CREATE INDEX IF NOT EXISTS idx_operators_email ON operators(email);
CREATE INDEX IF NOT EXISTS idx_operators_active ON operators(active);

-- Add comments for documentation
COMMENT ON TABLE operators IS 'Operators who can authenticate and manage persons in the system';
COMMENT ON COLUMN operators.username IS 'Unique username for login';
COMMENT ON COLUMN operators.email IS 'Unique email address';
COMMENT ON COLUMN operators.password IS 'Bcrypt hashed password (cost: 10)';
COMMENT ON COLUMN operators.active IS 'Whether the operator account is active';
