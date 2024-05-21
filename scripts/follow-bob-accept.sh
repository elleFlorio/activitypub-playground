#!/bin/sh

id=$1

# Accept a follow request sent to Bob. The id of the request should be specified
curl -v --data '{
  "@context": "https://www.w3.org/ns/activitystreams",
  "type": "Accept",
  "object": "'"$id"'",
  "actor": "http://anothercooldomain.com:8080/users/bobby"
}' \
localhost:8081/users/bobby/outbox
