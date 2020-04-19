package apicall

import (
	"log"

	"google.golang.org/api/gmail/v1"
)

func Label(name string, svc *gmail.Service) *gmail.Label {
	label, err := svc.Users.Labels.Get("me", name).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve label:\n %v", err)
	}
	return label
}
