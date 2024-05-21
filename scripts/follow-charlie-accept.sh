#!/bin/sh

id=$1

# Accept a follow request sent to Bob. The id of the request should be specified
curl -v --data '{
  "@context": "https://www.w3.org/ns/activitystreams",
  "type": "Accept",
  "object": "'"$id"'",
  "actor": "http://cooldomain.com:8080/users/chrlz"
}' \
localhost:8080/users/chrlz/outbox
