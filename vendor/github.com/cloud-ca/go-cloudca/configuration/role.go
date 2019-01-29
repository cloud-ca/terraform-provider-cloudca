package configuration

type Role struct {
	Id           string       `json:"id,omitempty"`
	Name         string       `json:"name,omitempty"`
	Environment  Environment  `json:"environment,omitempty"`
	Users        []User       `json:"users"`
	Organization Organization `json:"organization,omitempty"`
}
