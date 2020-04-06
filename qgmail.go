package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

var redirectURL string

// Request a token from the web, and return the received token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {

	// Open user-consent page in browser
	authURL := config.AuthCodeURL("Must be a random string", oauth2.AccessTypeOffline)
	openURL(authURL)

	// Fetch the authorization code
	var authCode string
	fmt.Println("Paste the authorization code:")
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	token, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return token
}

// Open URL using default browser
func openURL(url string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatalf("Unable to open browser: %v\n", err)
		fmt.Printf("Go to the following link in your browser then type the "+
			"authorization code: \n%v\n", url)
	}
}

// Retrieve token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	token := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(token)
	return token, err
}

// Save token to a file path.
func saveToken(path string, token *oauth2.Token) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func main() {
	var redirectPort int = 5000
	redirectURL = "http://localhost:" + strconv.Itoa(redirectPort)

	// Load app credentials.
	credentials, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// Create config struct.
	config, err := google.ConfigFromJSON(credentials, gmail.GmailLabelsScope)
	// config.RedirectURL = "http://localhost:" + strconv.Itoa(redirectPort)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	// Retrieve token.
	tokenFile := "token.json"
	token, err := tokenFromFile(tokenFile)
	if err != nil {
		token = getTokenFromWeb(config)
		saveToken(tokenFile, token)
	}

	// Create a service.
	ctx := context.TODO()
	tokenSource := config.TokenSource(context.TODO(), token)
	service, err := gmail.NewService(ctx, option.WithTokenSource(tokenSource))
	if err != nil {
		log.Fatalf("Unable to create Gmail service: %v", err)
	}

	// Display unread count.
	label, err := service.Users.Labels.Get("me", "INBOX").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve label: %v", err)
	}
	fmt.Println(label.ThreadsUnread)
}
