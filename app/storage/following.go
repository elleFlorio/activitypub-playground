package storage

var following = make(map[string][]string, 100)

func AddToFollowing(username string, actorId string) {
	if val, ok := following[username]; ok {
		following[username] = append(val, actorId)
		return
	}

	following[username] = []string{actorId}
}

func GetFollowing(username string) []string {
	if val, ok := following[username]; ok {
		return val
	}

	return []string{}
}
