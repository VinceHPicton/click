CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE SCHEMA IF NOT EXISTS app;

CREATE TABLE app.auth_attempts (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    mobile VARCHAR(20) NOT NULL,
    one_time_code INTEGER NOT NULL DEFAULT (FLOOR(RANDOM() * 900000 + 100000))::INT,
    used_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE app.users (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    name VARCHAR(50),
    bio VARCHAR(10000),
    birth_date DATE,
    last_location_long DECIMAL(9,6),
    last_location_lat DECIMAL(9,6),
    mobile VARCHAR(20) NOT NULL UNIQUE,
    last_active TIMESTAMPTZ,
    email VARCHAR(255),
    sex SMALLINT,
    interested_in SMALLINT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX users_phone_unique_active
ON app.users(mobile)
WHERE deleted_at IS NULL;

CREATE TABLE app.refresh_tokens (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_id UUID REFERENCES app.users NOT NULL,
    token_hash TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ
);

CREATE TABLE app.matches (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_1_id UUID REFERENCES app.users,
    user_2_id UUID REFERENCES app.users,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE app.likes (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_id UUID REFERENCES app.users,
    target_id UUID REFERENCES app.users,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE app.user_photos (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_id UUID REFERENCES app.users,
    url VARCHAR(4000),
    hash VARCHAR(64) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE app.messages (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_id UUID REFERENCES app.users,
    target_id UUID REFERENCES app.users,
    message VARCHAR(4000),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);