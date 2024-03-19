package storage

var timelines = make(map[string][]string, 100)

func AddToTimeline(username string, objectId string) {
	if val, ok := timelines[username]; ok {
		timelines[username] = append(val, objectId)
		return
	}

	timelines[username] = []string{objectId}
}

func GetTimeline(username string) ([]string, bool) {
	timeline, ok := timelines[username]
	return timeline, ok
}
