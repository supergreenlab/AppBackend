#!/bin/bash

docker run --rm -it --network cba45d93006e3385cc3ee52a8f68eb4535a89668030679c37ac2f4f8912951fe --env-file ../SuperGreenCloud/env.development --name minio -v supergreencloud_minio_data:/data minio/minio server /data
