package main

import (
	"log"

	"google.golang.org/api/gmail/v1"
)

func getLabel(name string, service *gmail.Service) *gmail.Label {
	label, err := service.Users.Labels.Get("me", name).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve label: %v", err)
	}
	return label
}
