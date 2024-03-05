package storage

import "elleFlorio/activitypub-playground/app/model"

var activities = make(map[string]model.Activity, 100)

func AddActivity(activity model.Activity) {
	activities[activity.Id] = activity
}

func GetActivity(id string) (model.Activity, bool) {
	activity, ok := activities[id]
	return activity, ok
}
