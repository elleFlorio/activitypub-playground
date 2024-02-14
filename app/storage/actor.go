package storage

import (
	"elleFlorio/activitypub-playground/app/model"
)

var actors = make(map[string]model.Actor, 100)

func CreateActor(user model.Actor) {
	actors[user.Id] = user
}

func GetActor(id string) (model.Actor, bool) {
	actor, ok := actors[id]
	return actor, ok
}
