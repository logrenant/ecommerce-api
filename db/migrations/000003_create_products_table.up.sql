CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), 
    name TEXT NOT NULL,
    description TEXT,
    price FLOAT NOT NULL,
    image_url VARCHAR,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT NOW()
);