#!/bin/bash

docker run --rm -it --network supergreencloud_back-tier --entrypoint=/bin/sh minio/mc
