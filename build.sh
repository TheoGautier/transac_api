#!/bin/bash

function build() {
  docker build -t "transac_api":latest -f ./Dockerfile .
}

function run() {
  docker-compose up
}
