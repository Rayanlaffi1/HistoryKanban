#!/bin/sh
psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -c "CREATE DATABASE keycloak"
