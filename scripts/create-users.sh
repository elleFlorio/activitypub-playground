#!/bin/sh

# Creates user Alice on the cooldomain.com server 
curl -v --data '{ "Name": "Alice", "Username": "alix" }' \
--header 'Content-Type: application/json' \
localhost:8080/users

# Creates user Bob on the anothercooldomain.com server
curl -v --data '{ "Name": "Bob", "Username": "bobby" }' \
--header 'Content-Type: application/json' \
localhost:8081/users

# Creates user Charlie on the cooldomain.com server
curl -v --data '{ "Name": "Charlie", "Username": "chrlz" }' \
--header 'Content-Type: application/json' \
localhost:8080/users
