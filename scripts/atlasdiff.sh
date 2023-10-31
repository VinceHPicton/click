#!/bin/bash

. ./config/env.sh

atlas migrate diff $1 --dev-url "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" --dir "file://db/migrations" --to "file://db/schema/schema.sql"