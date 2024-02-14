#!/bin/sh

docker-compose build --no-cache && docker-compose up && docker-compose down
