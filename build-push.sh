#!/bin/bash

# login to push images to Docker Hub, requires entering username and password
docker login

# build and push ledger account api docker image
docker build \
  --tag micrommath/ledger-account-api:latest \
  --file ./account-service/Dockerfile \
  ./account-service
docker image push micrommath/ledger-account-api:latest


# build and push ledger api docker image
docker build \
  --tag micrommath/ledger-ledger-api:latest \
  --file ./ledger-service/Dockerfile \
  ./ledger-service
docker image push micrommath/ledger-ledger-api:latest
