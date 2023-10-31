-- Create "likes" table
CREATE TABLE "likes" (
  "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
  "user_id" bigint NOT NULL,
  "target_id" bigint NOT NULL,
  "created_at" timestamp NOT NULL,
  PRIMARY KEY ("id")
);
-- Create "matches" table
CREATE TABLE "matches" (
  "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
  "user_1_id" bigint NOT NULL,
  "user_2_id" bigint NOT NULL,
  "created_at" timestamp NOT NULL,
  PRIMARY KEY ("id")
);
-- Create "messages" table
CREATE TABLE "messages" (
  "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
  "user_id" bigint NOT NULL,
  "target_id" bigint NOT NULL,
  "message" character varying(4000) NULL,
  "created_at" timestamp NOT NULL,
  PRIMARY KEY ("id")
);
-- Create "user_photos" table
CREATE TABLE "user_photos" (
  "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
  "user_id" bigint NOT NULL,
  "url" character varying(4000) NULL,
  "hash" character varying(64) NOT NULL,
  "created_at" timestamp NOT NULL,
  PRIMARY KEY ("id")
);
-- Create "users" table
CREATE TABLE "users" (
  "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
  "name" character varying(50) NOT NULL,
  "bio" character varying(1000) NOT NULL,
  "birth_date" date NOT NULL,
  "last_location_long" numeric(9,6) NULL,
  "last_location_lat" numeric(9,6) NULL,
  "mobile" character varying(20) NOT NULL,
  "last_active" timestamp NULL,
  "email" character varying(100) NOT NULL,
  "sex" smallint NOT NULL,
  "interested_in" smallint NOT NULL,
  "created_at" timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id")
);
