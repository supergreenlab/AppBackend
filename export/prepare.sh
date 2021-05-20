#!/bin/bash

docker-compose up .
docker run -i --rm --net=export_back-tier -e PGPASSWORD=password postgres psql -h postgres -U postgres sglapp < db_dump.sql
