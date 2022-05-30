package model

type Location struct {
	UUID     string `json:"uuid"`
	Name     string `json:"name"`
	Region   string `json:"region"`
	District string `json:"district"`
}
