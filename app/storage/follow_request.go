package storage

import "elleFlorio/activitypub-playground/app/model"

var followRequests = make(map[string]model.Activity, 100)

func AddFollowRequest(activity model.Activity) {
	followRequests[activity.Id] = activity
}

func GetFollowRequest(id string) (model.Activity, bool) {
	activity, ok := followRequests[id]
	return activity, ok
}

func GetFollowRequestsByActor(actorId string) []model.Activity {
	requests := make([]model.Activity, 0, 100)
	for _, request := range followRequests {
		if request.Object == actorId {
			requests = append(requests, request)
		}
	}

	return requests
}

func DeleteFollowRequest(id string) {
	delete(followRequests, id)
}
