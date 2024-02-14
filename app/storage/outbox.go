package storage

import (
	"elleFlorio/activitypub-playground/app/model"
)

var outbox = make(map[string][]model.Activity, 100)

func AddToOutbox(username string, activity model.Activity) {
	if val, ok := outbox[username]; ok {
		outbox[username] = append(val, activity)
		return
	}

	outbox[username] = []model.Activity{activity}
}

func GetOutbox(username string) []model.Activity {
	if val, ok := outbox[username]; ok {
		return val
	}

	return []model.Activity{}
}
