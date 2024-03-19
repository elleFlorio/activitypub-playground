package storage

import (
	"elleFlorio/activitypub-playground/app/model"
)

var objects = make(map[string]model.Object, 100)

func AddObject(object model.Object) {
	objects[object.Id] = object
}

func GetObject(id string) (model.Object, bool) {
	object, ok := objects[id]
	return object, ok
}

func GetObjectsByActorIdAndType(actorId string, objType string) []model.Object {
	objectsByActorAndType := make([]model.Object, 0, 100)
	for _, object := range objects {
		if object.AttributedTo == actorId && object.Type == objType {
			objectsByActorAndType = append(objectsByActorAndType, object)
		}
	}

	return objectsByActorAndType
}
