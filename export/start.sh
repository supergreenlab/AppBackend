#!/bin/bash

if [ "$#" -ne 1 ]; then
  echo "Usage: $0 /minio"
  exit
fi

docker run  --name=export --network=export_back-tier -p 8080:8080 --rm -it -v $(pwd)/config:/etc/export -v $1:/minio -v $(pwd)/sgl_export:/sgl_export supergreenlab/export
