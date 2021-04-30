#!/bin/bash

docker build -t timelapse-dev . -f Dockerfile.timelapse.dev
docker run  --name=timelapse -p 8083:8083 --rm -it -v $(pwd)/config:/etc/timelapse -v $(pwd):/app -v $(pwd)/timelapse:/var/timelapse timelapse-dev
docker rmi timelapse-dev
