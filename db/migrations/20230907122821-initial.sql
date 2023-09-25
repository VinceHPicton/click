
-- +migrate Up
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    bio VARCHAR(1000),
    birth_date DATE NOT NULL,
    last_location_long DECIMAL(9,6),
    last_location_lat DECIMAL(9,6),
    mobile VARCHAR(20),
    last_active TIMESTAMP,
    email VARCHAR(100),
    sex BIT NOT NULL,
    interested_in SMALLINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE matches (
    id BIGINT PRIMARY KEY,
    user_1_id BIGINT NOT NULL,
    user_2_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL
);

CREATE TABLE likes (
    id BIGINT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    target_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL
);

CREATE TABLE user_photos (
    id BIGINT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    url VARCHAR(4000),
    hash VARCHAR(64) NOT NULL,
    created_at TIMESTAMP NOT NULL
);

CREATE TABLE messages (
    id BIGINT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    target_id BIGINT NOT NULL,
    message VARCHAR(4000),
    created_at TIMESTAMP NOT NULL
);
-- +migrate Down
DROP TABLE users;

DROP TABLE matches;

DROP TABLE likes;

DROP TABLE user_photos;

DROP TABLE messages;