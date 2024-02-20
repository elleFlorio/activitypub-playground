package app

import (
	"elleFlorio/activitypub-playground/app/model"
	"elleFlorio/activitypub-playground/app/storage"
	"elleFlorio/activitypub-playground/config"
	"fmt"
	"log"

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

func newFollowActivity(username string, objectId string) (model.Activity, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		log.Fatal("Error generating Follow activity ID")
		return model.Activity{}, err
	}
	actorId := getActorId(username)

	followActivity := model.Activity{
		Id:     id.String(),
		Type:   "Follow",
		Actor:  actorId,
		Object: objectId,
	}

	return followActivity, nil
}

func AddToOutbox(username string, activity model.Activity) {
	id, _ := uuid.NewRandom()
	activity.Id = id.String()
	storage.AddToOutbox(username, activity)

}

func processActivity(activity model.Activity) {
	switch activity.Type {
	case "Follow":

	}
}

func sendFollowActivity(activity model.Activity) {
	isLocalUser := isLocal(activity.Object)
	if isLocalUser {
		log.Default().Println("Local user")
	}

}
