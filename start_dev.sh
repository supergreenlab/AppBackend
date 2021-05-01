#!/bin/bash

docker build -t appbackend-dev . -f Dockerfile.dev
docker run  --name=appbackend --network=supergreencloud_back-tier -p 8080:8080 --rm -it -v $(pwd)/config:/etc/appbackend -v $(pwd):/app appbackend-dev
