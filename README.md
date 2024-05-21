# ActivityPub Playground

This is the implementation of an [ActivityPub](https://www.w3.org/TR/activitypub/) server mainly used to understand how the protocol works. This is an educational prototype, so there are a lot compromises on how it works. Nevertheless, it can be a good toy project to understand the ActivityPub protocol for decentralized social networks.

This quick guide is intended to describe how to run the project and what it does with some practical examples. I will write an article on [my blog](florio.dev) to explain the implementation details.

## How to run
The only requirement to run the project is [Docker](https://www.docker.com/). Once you have Docker up and running simply run the `run.sh` script in the root folder of the repository.
```shell
./run.sh
```

## What it does
The project will run 2 ActivityPub servers written from scratch in [Go](https://go.dev/). The servers will run on 2 different domains inside the network created by Docker. Every server is mapped on the localhost to different ports:
- `cooldomain.com` --> `localhost:8080`
- `anothercooldomain.com` --> `localhost:8081`

Once the servers are up and running, it is possible to perform several actions hitting dedicated REST APIs. The actions available are the following:
- Create a user
- Search for a user using the [WebFinger](https://webfinger.net/) protocol
- Send a follow request to a user
- Accept a follow request
- Get the followers and following lists
- Create a post
- Check the timeline of a user. The timeline is a list of post created by users followed by the current one

The most notable thing that is missing is authentication, but I'll leave it for future development (maybe).

## Example workflow
In the `/scripts` folder there are several scripts that can be used as a reference on how to perform the actions mentioned above. Let's try them.

### Create users
Let's create three users: Alice and Charlie in the `cooldomain.com` server, and Bob in the `anothercooldomain.com` server.  
**Script**: `create-users.sh`
```shell
./create-users.sh
*   Trying 127.0.0.1:8080...
* Connected to localhost (127.0.0.1) port 8080 (#0)
> POST /users HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/8.1.2
> Accept: */*
> Content-Type: application/json
> Content-Length: 39
>
< HTTP/1.1 201 Created
< Content-Type: application/json; charset=utf-8
< Date: Tue, 21 May 2024 19:50:45 GMT
< Content-Length: 411
<
{
    "@context": "https://www.w3.org/ns/activitystreams",
    "id": "http://cooldomain.com:8080/users/alix",
    "type": "Person",
    "inbox": "http://cooldomain.com:8080/users/alix/inbox",
    "outbox": "http://cooldomain.com:8080/users/alix/outbox",
    "following": "http://cooldomain.com:8080/users/alix/following",
    "followers": "http://cooldomain.com:8080/users/alix/followers",
    "name": "Alice"
* Connection #0 to host localhost left intact
}*   Trying 127.0.0.1:8081...
* Connected to localhost (127.0.0.1) port 8081 (#0)
> POST /users HTTP/1.1
> Host: localhost:8081
> User-Agent: curl/8.1.2
> Accept: */*
> Content-Type: application/json
> Content-Length: 38
>
< HTTP/1.1 201 Created
< Content-Type: application/json; charset=utf-8
< Date: Tue, 21 May 2024 19:50:45 GMT
< Content-Length: 449
<
{
    "@context": "https://www.w3.org/ns/activitystreams",
    "id": "http://anothercooldomain.com:8080/users/bobby",
    "type": "Person",
    "inbox": "http://anothercooldomain.com:8080/users/bobby/inbox",
    "outbox": "http://anothercooldomain.com:8080/users/bobby/outbox",
    "following": "http://anothercooldomain.com:8080/users/bobby/following",
    "followers": "http://anothercooldomain.com:8080/users/bobby/followers",
    "name": "Bob"
* Connection #0 to host localhost left intact
}*   Trying 127.0.0.1:8080...
* Connected to localhost (127.0.0.1) port 8080 (#0)
> POST /users HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/8.1.2
> Accept: */*
> Content-Type: application/json
> Content-Length: 42
>
< HTTP/1.1 201 Created
< Content-Type: application/json; charset=utf-8
< Date: Tue, 21 May 2024 19:50:45 GMT
< Content-Length: 418
<
{
    "@context": "https://www.w3.org/ns/activitystreams",
    "id": "http://cooldomain.com:8080/users/chrlz",
    "type": "Person",
    "inbox": "http://cooldomain.com:8080/users/chrlz/inbox",
    "outbox": "http://cooldomain.com:8080/users/chrlz/outbox",
    "following": "http://cooldomain.com:8080/users/chrlz/following",
    "followers": "http://cooldomain.com:8080/users/chrlz/followers",
    "name": "Charlie"
* Connection #0 to host localhost left intact
}%
```
The three users have been created along with all the ActivityPub collections such as Inbox, Outbox, etc.

### Search for a user
We can search a user through the WebFinger protocol. To make things interesting, we are going to query the `cooldomain.com` server for Bob's user that is living in the `anothercooldomain.com` server.  
**Script**: `webfinger.sh`
```shell
./webfinger.sh
{
    "subject": "acct:bobby@anothercooldomain.com",
    "aliases": [
        "http://anothercooldomain.com:8080/users/bobby"
    ],
    "properties": null,
    "links": [
        {
            "rel": "self",
            "type": "application/activity+json",
            "href": "http://anothercooldomain.com:8080/users/bobby"
        }
    ]
}%
```
Looking at the Docker logs, we can see that the `cooldomain.com` server can resolve the place where Bob's user is living and query it to return the information requested:
```shell
ap-server-anothercool-1  | [GIN] 2024/05/21 - 19:55:05 | 200 |     135.771µs |      172.18.0.2 | GET      "/.well-known/webfinger?resource=acct:bobby@anothercooldomain.com"
ap-server-cool-1         | [GIN] 2024/05/21 - 19:55:05 | 200 |     855.617µs |    192.168.65.1 | GET      "/.well-known/webfinger?resource=acct:bobby@anothercooldomain.com"
```

### Follow users
Time to get social. We are sending out few follow requests:
- Alice --> Charlie
- Alice --> Bob
- Charlie --> Bob

Bob is our superstar.  
**Scripts**: `follow-alice-charlie.sh`, `follow-alice-bob.sh`, `follow-charlie-bob.sh`
```shell
./follow-alice-charlie.sh
*   Trying 127.0.0.1:8080...
* Connected to localhost (127.0.0.1) port 8080 (#0)
> POST /users/alix/outbox HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/8.1.2
> Accept: */*
> Content-Length: 183
> Content-Type: application/x-www-form-urlencoded
>
< HTTP/1.1 202 Accepted
< Content-Type: application/json; charset=utf-8
< Location: http://cooldomain.com:8080/users/alix/activity/0cf6c96b-bb11-4a22-aa3b-a0356ffd5086
< Date: Tue, 21 May 2024 19:51:03 GMT
< Content-Length: 4
<
* Connection #0 to host localhost left intact
null%

./follow-alice-bob.sh
*   Trying 127.0.0.1:8080...
* Connected to localhost (127.0.0.1) port 8080 (#0)
> POST /users/alix/outbox HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/8.1.2
> Accept: */*
> Content-Length: 190
> Content-Type: application/x-www-form-urlencoded
>
< HTTP/1.1 201 Created
< Content-Type: application/json; charset=utf-8
< Location: http://cooldomain.com:8080/users/alix/activity/d342143b-e18b-46bb-8a81-943f5361c8f4
< Date: Tue, 21 May 2024 19:51:10 GMT
< Content-Length: 4
<
* Connection #0 to host localhost left intact
null%

./follow-charlie-bob.sh
Note: Unnecessary use of -X or --request, POST is already inferred.
*   Trying 127.0.0.1:8080...
* Connected to localhost (127.0.0.1) port 8080 (#0)
> POST /users/chrlz/outbox HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/8.1.2
> Accept: */*
> Content-Length: 191
> Content-Type: application/x-www-form-urlencoded
>
< HTTP/1.1 201 Created
< Content-Type: application/json; charset=utf-8
< Location: http://cooldomain.com:8080/users/chrlz/activity/c690ef64-3824-4c49-9f86-5b360341e587
< Date: Tue, 21 May 2024 19:51:16 GMT
< Content-Length: 4
<
* Connection #0 to host localhost left intact
null%
```
Cool, we are social. But nothing matters if the requests are not accepted...

### Check and accept follow requests
Now we are going to check the follow requests of Charlie and Bob and accept them.  
**Scripts**: `follow-charlie-requests.sh`, `follow-charlie-accept.sh $requestID`, `follow-bob-requests.sh`, `follow-bob-accept.sh $requestID`
```shell
./follow-charlie-requests.sh
*   Trying 127.0.0.1:8080...
* Connected to localhost (127.0.0.1) port 8080 (#0)
> GET /users/chrlz/followers/requests HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/8.1.2
> Accept: */*
>
< HTTP/1.1 200 OK
< Content-Type: application/json; charset=utf-8
< Date: Tue, 21 May 2024 19:51:51 GMT
< Content-Length: 427
<
[
    {
        "@context": "https://www.w3.org/ns/activitystreams",
        "Id": "http://cooldomain.com:8080/users/alix/activity/0cf6c96b-bb11-4a22-aa3b-a0356ffd5086",
        "Type": "Follow",
        "Actor": "http://cooldomain.com:8080/users/alix",
        "Object": "http://cooldomain.com:8080/users/chrlz",
        "Target": "",
        "To": null,
        "Cc": null,
        "Published": "0001-01-01T00:00:00Z"
    }
* Connection #0 to host localhost left intact
]%

./follow-charlie-accept.sh http://cooldomain.com:8080/users/alix/activity/0cf6c96b-bb11-4a22-aa3b-a0356ffd5086
*   Trying 127.0.0.1:8080...
* Connected to localhost (127.0.0.1) port 8080 (#0)
> POST /users/chrlz/outbox HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/8.1.2
> Accept: */*
> Content-Length: 229
> Content-Type: application/x-www-form-urlencoded
>
< HTTP/1.1 202 Accepted
< Content-Type: application/json; charset=utf-8
< Location: http://cooldomain.com:8080/users/chrlz/activity/0139dcab-d3e8-4f38-9ad6-4138605521e9
< Date: Tue, 21 May 2024 19:52:06 GMT
< Content-Length: 4
<
* Connection #0 to host localhost left intact
null%

./follow-bob-requests.sh
*   Trying 127.0.0.1:8081...
* Connected to localhost (127.0.0.1) port 8081 (#0)
> GET /users/bobby/followers/requests HTTP/1.1
> Host: localhost:8081
> User-Agent: curl/8.1.2
> Accept: */*
>
< HTTP/1.1 200 OK
< Content-Type: application/json; charset=utf-8
< Date: Tue, 21 May 2024 19:52:34 GMT
< Content-Length: 868
<
[
    {
        "@context": "https://www.w3.org/ns/activitystreams",
        "Id": "http://cooldomain.com:8080/users/chrlz/activity/c690ef64-3824-4c49-9f86-5b360341e587",
        "Type": "Follow",
        "Actor": "http://cooldomain.com:8080/users/chrlz",
        "Object": "http://anothercooldomain.com:8080/users/bobby",
        "Target": "",
        "To": null,
        "Cc": null,
        "Published": "0001-01-01T00:00:00Z"
    },
    {
        "@context": "https://www.w3.org/ns/activitystreams",
        "Id": "http://cooldomain.com:8080/users/alix/activity/d342143b-e18b-46bb-8a81-943f5361c8f4",
        "Type": "Follow",
        "Actor": "http://cooldomain.com:8080/users/alix",
        "Object": "http://anothercooldomain.com:8080/users/bobby",
        "Target": "",
        "To": null,
        "Cc": null,
        "Published": "0001-01-01T00:00:00Z"
    }
* Connection #0 to host localhost left intact
]%

./follow-bob-accept.sh http://cooldomain.com:8080/users/chrlz/activity/c690ef64-3824-4c49-9f86-5b360341e587
*   Trying 127.0.0.1:8081...
* Connected to localhost (127.0.0.1) port 8081 (#0)
> POST /users/bobby/outbox HTTP/1.1
> Host: localhost:8081
> User-Agent: curl/8.1.2
> Accept: */*
> Content-Length: 237
> Content-Type: application/x-www-form-urlencoded
>
< HTTP/1.1 202 Accepted
< Content-Type: application/json; charset=utf-8
< Location: http://anothercooldomain.com:8080/users/bobby/activity/3dbbfd25-50e0-4d13-8091-7f5370141f5e
< Date: Tue, 21 May 2024 19:52:53 GMT
< Content-Length: 4
<
* Connection #0 to host localhost left intact
null%

./follow-bob-accept.sh http://cooldomain.com:8080/users/alix/activity/d342143b-e18b-46bb-8a81-943f5361c8f4
*   Trying 127.0.0.1:8081...
* Connected to localhost (127.0.0.1) port 8081 (#0)
> POST /users/bobby/outbox HTTP/1.1
> Host: localhost:8081
> User-Agent: curl/8.1.2
> Accept: */*
> Content-Length: 236
> Content-Type: application/x-www-form-urlencoded
>
< HTTP/1.1 202 Accepted
< Content-Type: application/json; charset=utf-8
< Location: http://anothercooldomain.com:8080/users/bobby/activity/82a261e9-dc78-4e1f-b252-fc2a9563c85d
< Date: Tue, 21 May 2024 19:53:16 GMT
< Content-Length: 4
<
* Connection #0 to host localhost left intact
null%
``` 
As you may have noticed, to accept a request we need the ID of the activity referencing the actual follow request. We accepted everything and now Bob is almost an influencer!

### Check following and followers
Let's verify the followers of every user. In case of Alice, since she has no followers, we are going to check who she is following.  
**Scripts**: `follow-alice-following.sh`, `follow-charlie-followers.sh`, `follow-bob-followers.sh`
```shell
./follow-alice-following.sh
*   Trying 127.0.0.1:8080...
* Connected to localhost (127.0.0.1) port 8080 (#0)
> GET /users/alix/following HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/8.1.2
> Accept: */*
>
< HTTP/1.1 200 OK
< Content-Type: application/json; charset=utf-8
< Date: Tue, 21 May 2024 19:53:33 GMT
< Content-Length: 101
<
[
    "http://cooldomain.com:8080/users/chrlz",
    "http://anothercooldomain.com:8080/users/bobby"
* Connection #0 to host localhost left intact
]%

./follow-charlie-followers.sh
*   Trying 127.0.0.1:8080...
* Connected to localhost (127.0.0.1) port 8080 (#0)
> GET /users/chrlz/followers HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/8.1.2
> Accept: */*
>
< HTTP/1.1 200 OK
< Content-Type: application/json; charset=utf-8
< Date: Tue, 21 May 2024 19:53:47 GMT
< Content-Length: 47
<
[
    "http://cooldomain.com:8080/users/alix"
* Connection #0 to host localhost left intact
]%

./follow-bob-followers.sh
*   Trying 127.0.0.1:8081...
* Connected to localhost (127.0.0.1) port 8081 (#0)
> GET /users/bobby/followers HTTP/1.1
> Host: localhost:8081
> User-Agent: curl/8.1.2
> Accept: */*
>
< HTTP/1.1 200 OK
< Content-Type: application/json; charset=utf-8
< Date: Tue, 21 May 2024 19:53:56 GMT
< Content-Length: 93
<
[
    "http://cooldomain.com:8080/users/chrlz",
    "http://cooldomain.com:8080/users/alix"
* Connection #0 to host localhost left intact
]%
```
Everything as planned.

### Making a post
What social network would this be without some posts? Let's make Charlie and Bob post something interesting.  
**Scripts**: `post-note-charlie.sh`, `post-note-bob.sh`
```shell
./post-note-charlie.sh
*   Trying 127.0.0.1:8080...
* Connected to localhost (127.0.0.1) port 8080 (#0)
> POST /users/chrlz/post HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/8.1.2
> Accept: */*
> Content-Length: 227
> Content-Type: application/x-www-form-urlencoded
>
< HTTP/1.1 201 Created
< Content-Type: application/json; charset=utf-8
< Location: http://cooldomain.com:8080/users/chrlz/activity/2e631dee-9c73-42eb-aa20-73e4ef4c41f6
< Date: Tue, 21 May 2024 19:54:06 GMT
< Content-Length: 4
<
* Connection #0 to host localhost left intact
null%

./post-note-bob.sh
*   Trying 127.0.0.1:8081...
* Connected to localhost (127.0.0.1) port 8081 (#0)
> POST /users/bobby/post HTTP/1.1
> Host: localhost:8081
> User-Agent: curl/8.1.2
> Accept: */*
> Content-Length: 237
> Content-Type: application/x-www-form-urlencoded
>
< HTTP/1.1 201 Created
< Content-Type: application/json; charset=utf-8
< Location: http://anothercooldomain.com:8080/users/bobby/activity/c08db357-316e-4882-8e27-31f0d61d24bd
< Date: Tue, 21 May 2024 19:54:26 GMT
< Content-Length: 4
<
* Connection #0 to host localhost left intact
null%
```
Some posts have been created. Time to check if they have been delivered to the followers..

### Get the user timeline
According to the followers we setup, this is what we expect:
- Alice will see both Charlie and Bob posts
- Charlie will see only Bob post

Let's verify this.  
**Scripts**: `timeline-alice.sh`, `timeline-charlie.sh`
```shell
./timeline-alice.sh
*   Trying 127.0.0.1:8080...
* Connected to localhost (127.0.0.1) port 8080 (#0)
> GET /users/alix/timeline HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/8.1.2
> Accept: */*
>
< HTTP/1.1 200 OK
< Content-Type: application/json; charset=utf-8
< Date: Tue, 21 May 2024 19:54:34 GMT
< Content-Length: 1146
<
[
    {
        "@context": "https://www.w3.org/ns/activitystreams",
        "Id": "http://cooldomain.com:8080/users/chrlz/object/624c311b-bfd2-43c4-a26a-cea739f5baa3",
        "Type": "Note",
        "Actor": "",
        "Object": "",
        "Target": "",
        "Name": "",
        "Content": "I think ActivityPub is super cool.",
        "Published": "2024-05-21T19:54:06Z",
        "AttributedTo": "http://cooldomain.com:8080/users/chrlz",
        "To": [
            "http://cooldomain.com:8080/users/chrlz/followers"
        ],
        "Cc": null
    },
    {
        "@context": "https://www.w3.org/ns/activitystreams",
        "Id": "http://anothercooldomain.com:8080/users/bobby/object/b94fdbcf-acd1-4074-b124-5f9b1f7a15a6",
        "Type": "Note",
        "Actor": "",
        "Object": "",
        "Target": "",
        "Name": "",
        "Content": "What a great vacation I had in Italy!",
        "Published": "2024-05-21T19:54:26Z",
        "AttributedTo": "http://anothercooldomain.com:8080/users/bobby",
        "To": [
            "http://anothercooldomain.com:8080/users/bobby/followers"
        ],
        "Cc": null
    }
* Connection #0 to host localhost left intact
]%

./timeline-charlie.sh
*   Trying 127.0.0.1:8080...
* Connected to localhost (127.0.0.1) port 8080 (#0)
> GET /users/chrlz/timeline HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/8.1.2
> Accept: */*
>
< HTTP/1.1 200 OK
< Content-Type: application/json; charset=utf-8
< Date: Tue, 21 May 2024 19:54:55 GMT
< Content-Length: 586
<
[
    {
        "@context": "https://www.w3.org/ns/activitystreams",
        "Id": "http://anothercooldomain.com:8080/users/bobby/object/b94fdbcf-acd1-4074-b124-5f9b1f7a15a6",
        "Type": "Note",
        "Actor": "",
        "Object": "",
        "Target": "",
        "Name": "",
        "Content": "What a great vacation I had in Italy!",
        "Published": "2024-05-21T19:54:26Z",
        "AttributedTo": "http://anothercooldomain.com:8080/users/bobby",
        "To": [
            "http://anothercooldomain.com:8080/users/bobby/followers"
        ],
        "Cc": null
    }
* Connection #0 to host localhost left intact
]%
```
Mic drop!!!
