package apicall

import (
	"google.golang.org/api/gmail/v1"
)

func Label(name string, svc *gmail.Service) (*gmail.Label, error) {
	label, err := svc.Users.Labels.Get("me", name).Do()
	return label, err
}
