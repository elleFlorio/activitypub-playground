package model

import (
	"encoding/json"
	"time"
)

type Actor struct {
	Id        string `json:"id"`
	Type      string `json:"type"`
	Inbox     string `json:"inbox"`
	Outbox    string `json:"outbox"`
	Following string `json:"following"`
	Followers string `json:"followers"`
	Name      string `json:"name"`
}

func (a *Actor) MarshalJSON() ([]byte, error) {
	type ActorAlias Actor
	return json.Marshal(&struct {
		Context string `json:"@context"`
		*ActorAlias
	}{
		Context:    "https://www.w3.org/ns/activitystreams",
		ActorAlias: (*ActorAlias)(a),
	})
}

type Activity struct {
	Id        string
	Type      string
	Actor     string
	Object    string
	Target    string
	To        []string
	Cc        []string
	Published time.Time
}

func (a *Activity) MarshalJSON() ([]byte, error) {
	type ActivityAlias Activity
	return json.Marshal(&struct {
		Context string `json:"@context"`
		*ActivityAlias
	}{
		Context:       "https://www.w3.org/ns/activitystreams",
		ActivityAlias: (*ActivityAlias)(a),
	})
}

type Object struct {
	Id           string
	Type         string
	Actor        string
	Object       string
	Target       string
	Name         string
	Content      string
	Published    time.Time
	AttributedTo string
	To           []string
	Cc           []string
}

func (a *Object) MarshalJSON() ([]byte, error) {
	type ObjectAlias Object
	return json.Marshal(&struct {
		Context string `json:"@context"`
		*ObjectAlias
	}{
		Context:     "https://www.w3.org/ns/activitystreams",
		ObjectAlias: (*ObjectAlias)(a),
	})
}
