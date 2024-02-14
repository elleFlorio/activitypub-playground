package app

import (
	"elleFlorio/activitypub-playground/app/model"
	"elleFlorio/activitypub-playground/config"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
)

const WebFingerEndpoint string = "/.well-known/webfinger"

func ResolveWellKnown(resource string) model.Account {
	username, domain := parseRerource(resource)
	if domain == config.Domain {
		return getAccount(username)
	}

	return findAccount(username, domain)
}

func parseRerource(resource string) (string, string) {
	pattern := regexp.MustCompile(`acct:(?P<username>[a-z]+)@(?P<domain>[a-z]+\.[a-z]+)`)
	matches := pattern.FindStringSubmatch(resource)

	username := matches[pattern.SubexpIndex("username")]
	domain := matches[pattern.SubexpIndex("domain")]

	return username, domain
}

func getAccount(username string) model.Account {
	actorRef := fmt.Sprintf("http://%s:8080/users/%s", config.Domain, username)
	subject := fmt.Sprintf("acct:%s@%s", username, config.Domain)
	aliases := []string{actorRef}
	links := []model.AccountRel{
		{
			Rel:  "self",
			Type: "application/activity+json",
			Href: actorRef,
		},
	}

	return model.Account{
		Subject: subject,
		Aliases: aliases,
		Links:   links,
	}
}

func findAccount(username string, domain string) model.Account {
	uri := fmt.Sprintf("http://%s:8080%s?resource=acct:%s@%s", domain, WebFingerEndpoint, username, domain)

	resp, err := http.Get(uri)
	if err != nil {
		log.Fatalf("Error resolving webfinger resource acct:%s@%s - %s", username, domain, err)
	}

	defer resp.Body.Close()

	var account model.Account

	if err := json.NewDecoder(resp.Body).Decode(&account); err != nil {
		log.Fatal("Error decoding webfinger response")
	}

	return account
}
