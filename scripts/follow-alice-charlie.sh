#!/bin/sh

# Alice issue a request to follow Charlie. They are on the same server
curl -v --data '{
  "@context": "https://www.w3.org/ns/activitystreams",
  "type": "Follow",
  "actor": "http://cooldomain.com:8080/users/alix",
  "object": "http://cooldomain.com:8080/users/chrlz"
}' \
localhost:8080/users/alix/outbox
