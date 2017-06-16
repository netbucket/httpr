#!/bin/bash
#
# Build the Docker image for httpr
#
HTTPR_IMAGE="netbucket/httpr"

docker build dist -f dist/Dockerfile -t $HTTPR_IMAGE:0.1.0 -t $HTTPR_IMAGE:latest
