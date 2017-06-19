#!/bin/bash
#
# Build the Docker image for httpr
#
HTTPR_IMAGE="netbucket/httpr"
HTTPR_VERSION="0.1.3"

docker build dist -f dist/Dockerfile -t $HTTPR_IMAGE:$HTTPR_VERSION -t $HTTPR_IMAGE:latest
docker push $HTTPR_IMAGE:$HTTPR_VERSION
docker push $HTTPR_IMAGE:latest
