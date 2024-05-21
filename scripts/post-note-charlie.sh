#!/bin/sh

published=$(date -u +%Y-%m-%dT%H:%M:%SZ)

# Create a post on Charlie profile
curl -v --data '{
  "@context": "https://www.w3.org/ns/activitystreams",
  "type": "Note",
  "content": "I think ActivityPub is super cool.",
  "published": "'"$published"'",
  "to": ["http://cooldomain.com:8080/users/chrlz/followers"]
}' \
localhost:8080/users/chrlz/post
