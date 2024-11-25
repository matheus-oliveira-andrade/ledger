#!/bin/bash

# login to push images to Docker Hub, requires entering username and password
docker login

# build and push transactions movement api docker image
docker build \
  --tag micrommath/ledger-account-api:latest \
  --file ./account-service/Dockerfile \
  ./account-service

docker image push micrommath/ledger-account-api:latest