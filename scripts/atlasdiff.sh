#!/bin/bash

source .env

# How to use this script:
# This script compares schema.sql (your desired database state) to what would be achieved by running 
# the migration files in /migrations consecutively.
# If the desired state does not equal this, atlas creates a new migration file named "[datetime]_[argument given to script ($1)].sql"
# You can edit this file if you want to, but you will need to rerun atlas migrate hash as the hashes in atlas.sum wont match if you edit.
# After this, you can run the migrate.sh script to apply changes

atlas migrate diff $1 \
--dir "file://db/migrations" \
--to "file://db/schema/0_schema.sql" \
--dev-url "docker://postgres/16"