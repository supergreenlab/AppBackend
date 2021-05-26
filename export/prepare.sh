#!/bin/bash

if [ "$#" -ne 1 ]; then
  echo "Usage: $0 pg_password"
  exit
fi

POSTGRES_PASSWORD="$1" docker-compose up -d
echo "create database sglapp;" | docker run -i --rm --net=export_back-tier -e PGPASSWORD="$1" postgres psql -h postgres -U postgres
docker run -i --rm --net=export_back-tier -e PGPASSWORD="$1" postgres psql -h postgres -U postgres sglapp < db_dump.sql
