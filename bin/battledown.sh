#!/bin/bash

PARENT_DIR=$(dirname $(dirname $(readlink -f "$0")))
PROJECT_NAME=$(basename "$PARENT_DIR")

echo ${PROJECT_NAME}

docker compose -p $PROJECT_NAME -f $PARENT_DIR/docker-compose.yml down