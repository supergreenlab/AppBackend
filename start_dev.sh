#!/bin/bash

docker build -t appbackend-dev . -f Dockerfile.dev
docker run -p 8080:8080 --rm -it -v $(pwd)/config:/etc/appbackend -v $(pwd):/app appbackend-dev
docker rmi appbackend-dev
