package app

import (
	"bytes"
	"elleFlorio/activitypub-playground/app/model"
	"elleFlorio/activitypub-playground/app/storage"
	"elleFlorio/activitypub-playground/config"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func getActorId(username string) string {
	return fmt.Sprintf("http://%s:8080/users/%s", config.Domain, username)
}

func newActor(username string) model.Actor {
	actorId := getActorId(username)
	return model.Actor{
		Id:        actorId,
		Type:      "Person",
		Inbox:     fmt.Sprintf("%s/inbox", actorId),
		Outbox:    fmt.Sprintf("%s/outbox", actorId),
		Following: fmt.Sprintf("%s/following", actorId),
		Followers: fmt.Sprintf("%s/followers", actorId),
	}
}

func AddToInbox(username string, activity model.Activity) int {
	storage.AddToInbox(username, activity)
	log.Default().Printf("Added to user %s inbox activity of type %s", username, activity.Type)

	return 201
}

func AddToOutbox(username string, activity model.Activity) (string, int) {
	id, _ := uuid.NewRandom()
	activity.Id = id.String()
	storage.AddToOutbox(username, activity)
	processActivity(activity)

	log.Default().Printf("Added to user %s outbox activity of type %s", username, activity.Type)

	return activity.Id, 201
}

func processActivity(activity model.Activity) {
	switch activity.Type {
	case "Follow":
		sendFollowActivity(activity)
	}
}

func sendFollowActivity(activity model.Activity) {
	isLocalUser := isLocal(activity.Object)
	if isLocalUser {
		username, _ := parseId(activity.Object)
		AddToInbox(username, activity)
	}

	toFollow, _ := getRemoteUser(activity.Object)

	body, _ := activity.MarshalJSON()
	_, err := http.Post(toFollow.Inbox, "application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\"", bytes.NewBuffer(body))
	if err != nil {
		log.Fatalf("Error sending to inbox: %s", toFollow.Inbox)
	}
}
