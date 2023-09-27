CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    bio VARCHAR(1000) NOT NULL,
    birth_date DATE NOT NULL,
    last_location_long DECIMAL(9,6),
    last_location_lat DECIMAL(9,6),
    mobile VARCHAR(20) NOT NULL,
    last_active TIMESTAMP,
    email VARCHAR(100) NOT NULL,
    sex BIT NOT NULL,
    interested_in SMALLINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE matches (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_1_id BIGINT NOT NULL,
    user_2_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL
);

CREATE TABLE likes (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_id BIGINT NOT NULL,
    target_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL
);

CREATE TABLE user_photos (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_id BIGINT NOT NULL,
    url VARCHAR(4000),
    hash VARCHAR(64) NOT NULL,
    created_at TIMESTAMP NOT NULL
);

CREATE TABLE messages (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_id BIGINT NOT NULL,
    target_id BIGINT NOT NULL,
    message VARCHAR(4000),
    created_at TIMESTAMP NOT NULL
);