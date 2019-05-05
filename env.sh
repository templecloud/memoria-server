#!/bin/bash

# See: go help modules
export GO111MODULE=on

export MONGO_NAME="mongodb"
export MONGO_VERSION="4.1.9-bionic"

function mongo-run() {
    docker run -p 27017:27017 --name "${MONGO_NAME}" -d mongo:"${MONGO_VERSION}"
}

function mongo-start() {
    docker start "${MONGO_NAME}"
}

function mongo-stop() {
    docker stop "${MONGO_NAME}"
}

function mongo-delete() {
    docker delete -f "${MONGO_NAME}"
}

function mongo-exec() {
    docker exec -it "${MONGO_NAME}" bash
}

function mongo-shell() {
    docker exec -it "${MONGO_NAME}" mongo
}
