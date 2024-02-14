package model

type Account struct {
	Subject    string       `json:"subject"`
	Aliases    []string     `json:"aliases"`
	Properties []string     `json:"properties"`
	Links      []AccountRel `json:"links"`
}

type AccountRel struct {
	Rel  string `json:"rel"`
	Type string `json:"type"`
	Href string `json:"href"`
}
