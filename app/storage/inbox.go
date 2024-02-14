package storage

import (
	"elleFlorio/activitypub-playground/app/model"
)

var inbox = make(map[string][]model.Activity, 100)

func AddToInbox(username string, activity model.Activity) {
	if val, ok := inbox[username]; ok {
		inbox[username] = append(val, activity)
		return
	}

	inbox[username] = []model.Activity{activity}
}

func GetInbox(username string) []model.Activity {
	if val, ok := inbox[username]; ok {
		return val
	}

	return []model.Activity{}
}
