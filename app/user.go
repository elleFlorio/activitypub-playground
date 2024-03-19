package app

import (
	"elleFlorio/activitypub-playground/app/model"
	"elleFlorio/activitypub-playground/app/storage"
	"elleFlorio/activitypub-playground/config"
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strings"
)

func CreateUser(username string, name string) model.Actor {
	actor := newActor(username)
	actor.Name = name
	storage.CreateActor(actor)
	return actor
}

func SearchUser(acct string) (model.Actor, int) {
	username, domain, ok := strings.Cut(acct, "@")
	if !ok {
		log.Fatalf("Bad account format: %s", acct)
		return model.Actor{}, http.StatusBadRequest
	}

	if domain == config.Domain {
		id := getActorId(username)
		if actor, ok := storage.GetActor(id); ok {
			return actor, http.StatusOK
		} else {
			log.Fatalf("User %s not found", id)
			return model.Actor{}, http.StatusNotFound
		}

	}

	account := findAccount(username, domain)
	for _, link := range account.Links {
		if link.Rel == "self" {
			return getRemoteUser(link.Href)
		}
	}

	return model.Actor{}, http.StatusNotFound

}

func getRemoteUser(id string) (model.Actor, int) {
	resp, err := http.Get(id)
	if err != nil {
		log.Fatalf("Error finding user: %s, error: %s", id, err)
		return model.Actor{}, http.StatusInternalServerError
	}

	if resp.StatusCode == http.StatusNotFound {
		log.Fatalf("User %s not found", id)
		return model.Actor{}, http.StatusNotFound
	}

	defer resp.Body.Close()

	var actor model.Actor

	if err := json.NewDecoder(resp.Body).Decode(&actor); err != nil {
		log.Fatal("Error decoding webfinger response")
		return model.Actor{}, http.StatusInternalServerError
	}

	return actor, http.StatusOK
}

func GetUser(username string) (model.Actor, int) {
	id := getActorId(username)
	if actor, ok := storage.GetActor(id); ok {
		return actor, http.StatusOK
	}

	log.Fatalf("User %s not found", id)
	return model.Actor{}, http.StatusNotFound
}

func isLocal(id string) bool {
	_, domain := parseId(id)
	return domain == config.Domain
}

func parseId(id string) (string, string) {
	pattern := regexp.MustCompile(`http://(?P<domain>[a-z]+\.[a-z]+):8080/users/(?P<username>[a-z]+)/*`)
	matches := pattern.FindStringSubmatch(id)

	username := matches[pattern.SubexpIndex("username")]
	domain := matches[pattern.SubexpIndex("domain")]

	return username, domain

}

func GetFollowRequests(username string) ([]model.Activity, int) {
	actorId := getActorId(username)
	followRequests := storage.GetFollowRequestsByActor(actorId)
	return followRequests, http.StatusOK
}

func GetFollowers(username string) ([]string, int) {
	return storage.GetFollowers(username), http.StatusOK
}

func GetFollowing(username string) ([]string, int) {
	return storage.GetFollowing(username), http.StatusOK
}

func GetPosts(username string) ([]model.Object, int) {
	actorId := getActorId(username)
	return storage.GetObjectsByActorIdAndType(actorId, "Note"), http.StatusOK
}

func GetTimeline(username string) ([]model.Object, int) {
	timeline := make([]model.Object, 0, 100)
	if noteIds, ok := storage.GetTimeline(username); ok {
		for _, noteId := range noteIds {
			if isLocal(noteId) {
				if note, ok := storage.GetObject(noteId); ok {
					timeline = append(timeline, note)
				} else {
					log.Default().Printf("WARNING: Cannot find Note object with id: %s\n", noteId)
				}
			} else {
				note, status := GetRemoteObject(noteId)
				if status != http.StatusOK {
					log.Fatalf("Error getting remote Note object with id: %s\n", noteId)
				}
				timeline = append(timeline, note)
			}
		}

		return timeline, http.StatusOK
	}

	return []model.Object{}, http.StatusNotFound
}
