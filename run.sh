#!/bin/bash

docker build -t orders-service .
docker compose up --build