package app

import "testing"

func TestParseResource(t *testing.T) {
	const testResource string = "acct:user@domain.test"
	const testUser string = "user"
	const testDomain string = "domain.test"

	username, domain := parseRerource(testResource)

	if username != testUser {
		t.Errorf("Error parsing username. Expected: %s, Actual: %s", testUser, username)
	}

	if domain != testDomain {
		t.Errorf("Error parsing domain. Expected: %s, Actual: %s", testDomain, domain)
	}
}

func TestResolveWellKnown(t *testing.T) {
	// TODO
}
