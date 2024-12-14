-- Users Table
CREATE TABLE users (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Signatures Table
CREATE TABLE signatures (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    template_data JSONB NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Links Table
CREATE TABLE links (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    signature_id UUID REFERENCES signatures(id),
    url TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Clicks Table
CREATE TABLE clicks (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    link_id UUID REFERENCES links(id),
    timestamp TIMESTAMP DEFAULT NOW(),
    ip_address TEXT
);
