package app

import (
	"fmt"
	"testing"
)

func TestParseId(t *testing.T) {
	const testUser string = "user"
	const testDomain string = "domain.test"
	testId := fmt.Sprintf("http://%s:8080/users/%s/something/somethingelse", testDomain, testUser)

	username, domain := parseId(testId)

	if username != testUser {
		t.Errorf("Error parsing username. Expected: %s, Actual: %s", testUser, username)
	}

	if domain != testDomain {
		t.Errorf("Error parsing domain. Expected: %s, Actual: %s", testDomain, domain)
	}
}
