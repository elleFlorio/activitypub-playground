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
	log.Default().Printf("Added to user %s inbox activity of type %s with id %s", username, activity.Type, activity.Id)

	return processInboxActivity(username, activity)
}

func processInboxActivity(username string, activity model.Activity) int {
	switch activity.Type {
	case "Follow":
		storage.AddFollowRequest(activity)
		log.Default().Printf("User %s received follow request from %s", username, activity.Actor)

		return http.StatusCreated
	case "Accept":
		log.Default().Printf("User %s accepted Follow request from %s", activity.Actor, username)

		return addToFollowing(username, activity)
	}

	return http.StatusMethodNotAllowed
}

func addToFollowing(username string, activity model.Activity) int {
	activities := storage.GetOutbox(username)
	for _, userActivity := range activities {
		if userActivity.Id == activity.Object {
			storage.AddToFollowing(username, activity.Actor)

			return http.StatusAccepted
		}
	}

	return http.StatusNotFound
}

func AddToOutbox(username string, activity model.Activity) (string, int) {
	id, _ := uuid.NewRandom()
	activity.Id = id.String()
	storage.AddToOutbox(username, activity)
	log.Default().Printf("Added to user %s outbox activity of type %s with id %s", username, activity.Type, activity.Id)

	result := processOutboxActivity(activity)

	return activity.Id, result
}

func processOutboxActivity(activity model.Activity) int {
	switch activity.Type {
	case "Follow":
		resp := sendToActor(activity.Object, activity)
		return resp
	case "Accept":
		username, _ := parseId(activity.Actor)
		if followActivity, ok := storage.GetFollowRequest(activity.Object); ok {
			resp := sendToActor(followActivity.Actor, activity)
			if resp == http.StatusAccepted {
				log.Default().Printf("User %s Accepted follow request", activity.Actor)
				storage.AddToFollowers(username, followActivity.Actor)
				storage.DeleteFollowRequest(followActivity.Id)
			}

			return resp
		}

		log.Default().Printf("Follow activity with id %s not found", activity.Object)

		return http.StatusNotFound
	case "Create":

	}

	return http.StatusMethodNotAllowed
}

func sendToActor(actorId string, activity model.Activity) int {
	isLocalUser := isLocal(actorId)
	if isLocalUser {
		username, _ := parseId(actorId)
		AddToInbox(username, activity)
	}

	toFollow, _ := getRemoteUser(actorId)

	body, _ := activity.MarshalJSON()
	resp, err := http.Post(toFollow.Inbox, "application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\"", bytes.NewBuffer(body))
	if err != nil {
		log.Fatalf("Error sending to inbox: %s", toFollow.Inbox)
	}

	return resp.StatusCode
}
