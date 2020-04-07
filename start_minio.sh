#!/bin/bash

docker run --rm -it -p 9000:9000 --network 589fd7238a2a92fb9668cdc917638581e19e6e284ffabac372ed41cc40bca47c --env-file ../SuperGreenCloud/env.development --name minio -v supergreencloud_minio_data:/data minio/minio server /data
