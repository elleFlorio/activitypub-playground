package storage

import "log"

var followers = make(map[string][]string, 100)

func AddToFollowers(username string, actorId string) {
	log.Default().Printf("Adding to %s followers actor %s\n", username, actorId)
	if val, ok := followers[username]; ok {
		followers[username] = append(val, actorId)
		return
	}

	followers[username] = []string{actorId}
}

func GetFollowers(username string) []string {
	if val, ok := followers[username]; ok {
		return val
	}

	return []string{}
}
