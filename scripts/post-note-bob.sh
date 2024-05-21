#!/bin/sh

published=$(date -u +%Y-%m-%dT%H:%M:%SZ)

# Create a post on Bob profile
curl -v --data '{
  "@context": "https://www.w3.org/ns/activitystreams",
  "type": "Note",
  "content": "What a great vacation I had in Italy!",
  "published": "'"$published"'",
  "to": ["http://anothercooldomain.com:8080/users/bobby/followers"]
}' \
localhost:8081/users/bobby/post
