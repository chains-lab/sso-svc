package models

type TokensPair struct {
	Refresh string `json:"refresh"`
	Access  string `json:"access"`
}
