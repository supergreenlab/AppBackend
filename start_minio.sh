#!/bin/bash

docker run --rm -it -p 9000:9000 --network supergreencloud_back-tier --env-file ../SuperGreenCloud/env.development --name minio -v supergreencloud_minio_data:/data minio/minio server /data
