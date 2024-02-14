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
	Id     string
	Type   string
	Actor  string
	Object string
	Target string
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
