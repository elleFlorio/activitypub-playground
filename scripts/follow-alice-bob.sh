#!/bin/sh

# Alice issue a request to follow Bob. They are on different servers
curl -v --data '{
  "@context": "https://www.w3.org/ns/activitystreams",
  "type": "Follow",
  "actor": "http://cooldomain.com:8080/users/alix",
  "object": "http://anothercooldomain.com:8080/users/bobby"
}' \
localhost:8080/users/alix/outbox
