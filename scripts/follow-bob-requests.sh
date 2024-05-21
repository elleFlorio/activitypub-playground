#!/bin/sh

# Get the list of follow requests sent to Bob
curl -v localhost:8081/users/bobby/followers/requests
