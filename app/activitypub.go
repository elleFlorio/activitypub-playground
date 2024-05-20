package app

import (
	"bytes"
	"elleFlorio/activitypub-playground/app/model"
	"elleFlorio/activitypub-playground/app/storage"
	"elleFlorio/activitypub-playground/config"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func getActorId(username string) string {
	return fmt.Sprintf("http://%s:8080/users/%s", config.Domain, username)
}

func getActivityId(username string) string {
	id, _ := uuid.NewRandom()
	return fmt.Sprintf("http://%s:8080/users/%s/activity/%s", config.Domain, username, id.String())
}

func getObjectId(username string) string {
	id, _ := uuid.NewRandom()
	return fmt.Sprintf("http://%s:8080/users/%s/object/%s", config.Domain, username, id.String())
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
	case "Create":
		log.Default().Printf("Added object %s to user %s timeline", activity.Object, username)
		storage.AddToTimeline(username, activity.Object)
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

func AddObjectToOutbox(username string, object model.Object) (string, int) {
	object.Id = getObjectId(username)
	object.AttributedTo = getActorId(username)
	processOutboxObject(object)

	activity := wrapInCreate(username, object)
	log.Default().Printf("Received for user %s object of type %s and wrapped in activity with id %s", username, object.Type, activity.Id)

	return AddActivityToOutbox(username, activity)
}

func processOutboxObject(object model.Object) {
	switch object.Type {
	case "Note":
		storage.AddObject(object)
	}
}

func AddActivityToOutbox(username string, activity model.Activity) (string, int) {
	activity.Id = getActivityId(username)
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
		// TODO simplification, we propagate only to followers
		followersId := activity.To[0]
		username, _ := parseId(followersId)
		followers := storage.GetFollowers(username)
		for _, follower := range followers {
			sendToActor(follower, activity)
		}
	}

	return http.StatusMethodNotAllowed
}

func sendToActor(actorId string, activity model.Activity) int {
	isLocalUser := isLocal(actorId)
	if isLocalUser {
		username, _ := parseId(actorId)
		AddToInbox(username, activity)
		return http.StatusAccepted
	}

	toFollow, _ := getRemoteUser(actorId)

	body, _ := activity.MarshalJSON()
	resp, err := http.Post(toFollow.Inbox, "application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\"", bytes.NewBuffer(body))
	if err != nil {
		log.Fatalf("Error sending to inbox: %s", toFollow.Inbox)
	}

	return resp.StatusCode
}

func wrapInCreate(username string, object model.Object) model.Activity {
	activityId := getActivityId(username)
	actorId := getActorId(username)

	createActivity := model.Activity{
		Id:        activityId,
		Type:      "Create",
		Actor:     actorId,
		Object:    object.Id,
		Published: object.Published,
		To:        object.To,
		Cc:        object.Cc,
	}

	return createActivity
}

func GetObjectById(username string, id string) (model.Object, int) {
	objectId := fmt.Sprintf("http://%s:8080/users/%s/object/%s", config.Domain, username, id)
	if object, ok := storage.GetObject(objectId); ok {
		return object, http.StatusOK
	}
	return model.Object{}, http.StatusNotFound
}

func GetRemoteObject(uri string) (model.Object, int) {
	resp, err := http.Get(uri)
	if err != nil {
		log.Fatalf("Error Getting remote object %s", uri)
	}

	defer resp.Body.Close()

	var remoteObject model.Object

	if err := json.NewDecoder(resp.Body).Decode(&remoteObject); err != nil {
		log.Fatal("Error decoding remote object response")
		return model.Object{}, http.StatusInternalServerError
	}

	return remoteObject, http.StatusOK
}
