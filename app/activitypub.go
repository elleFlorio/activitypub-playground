package app

import (
	"elleFlorio/activitypub-playground/app/model"
	"elleFlorio/activitypub-playground/config"
	"fmt"
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
