package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
	"time"

	"golang.org/x/oauth2"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type (
	authParams struct {
		CodeVerifier string
		Request      *authRequest
		Response     *authResponse
	}

	authRequest struct {
		State               string
		CodeChallenge       string
		CodeChallengeMethod string
		Scopes              []string
	}

	authResponse struct {
		State string
		Code  string
	}

	clientCreds struct {
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		AuthURL      string `json:"auth_uri"`
		TokenURL     string `json:"token_uri"`
	}
)

var (
	oauthConf *oauth2.Config
	mustPaste = false
)

func newAuthParams() *authParams {
	a := &authParams{
		CodeVerifier: generateRandomString(43),
		Request: &authRequest{
			State:               generateRandomString(10),
			CodeChallengeMethod: "S256",
			Scopes:              []string{gmail.GmailLabelsScope},
		},
		Response: &authResponse{},
	}
	a.Request.CodeChallenge = createCodeChallenge(a.CodeVerifier,
		a.Request.CodeChallengeMethod)
	return a
}

func newGmailService(config *oauth2.Config, token *oauth2.Token) (*gmail.Service, error) {
	tokenSource := config.TokenSource(context.TODO(), token)
	service, err := gmail.NewService(context.TODO(), option.WithTokenSource(tokenSource))
	if err != nil {
		log.Fatalf("Unable to create Gmail service: %v", err)
	}
	return service, err
}

func newOauthConf(credsFile string, conf *config) *oauth2.Config {
	var redirectURL string

	redirectPort, err := getRandomPort()
	// Use manual copy/paste method if error occurs.
	if err != nil {
		mustPaste = true
		redirectURL = "urn:ietf:wg:oauth:2.0:oob"
	} else {
		redirectURL = "http://localhost" + redirectPort
	}
	oauthConf := &oauth2.Config{
		RedirectURL: redirectURL,
		Scopes:      auth.Request.Scopes,
	}
	readCredentials(credsFile, oauthConf)
	return oauthConf
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	opts := []oauth2.AuthCodeOption{
		oauth2.SetAuthURLParam("access_type", "offline"),
		oauth2.SetAuthURLParam("code_challenge",
			auth.Request.CodeChallenge),
		oauth2.SetAuthURLParam("code_challenge_method",
			auth.Request.CodeChallengeMethod),
	}

	// Fetch the authorization code.
	authURL := config.AuthCodeURL(auth.Request.State, opts...)
	auth.Response.getAuthCode(authURL, mustPaste)

	// Verify state parameter
	if auth.Request.State != auth.Response.State {
		log.Fatal("Error: This request wasn't initialised by qGmail.")
	}

	token, err := config.Exchange(context.TODO(), auth.Response.Code,
		oauth2.SetAuthURLParam("code_verifier", auth.CodeVerifier))
	if err != nil {
		log.Fatalf("Error: Unable to retrieve token from the web.\n%v", err)
	}
	return token
}

func (a *authResponse) getAuthCode(authURL string, mustPaste bool) {
	if mustPaste {
		fmt.Println("Paste the authorization code here:")
		fmt.Scan(a.Code)
	} else {
		fmt.Println("Opening browser for user consent...\n" + authURL + "\n")
		openURL(authURL)
		time.Sleep(2 * time.Second)
		a.startHTTPListener()
	}
}

func openURL(url string) {
	switch runtime.GOOS {
	case "linux":
		exec.Command("xdg-open", url).Start()
	case "windows":
		exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		exec.Command("open", url).Start()
	default:
		fmt.Println("Error: Could not open your browser. Please open the" +
			" above URL manually and then proceed.")
	}
}

func readCredentials(credsFile string, conf *oauth2.Config) {
	var err error
	var j struct {
		Creds *clientCreds `json:"installed"`
	}

	if _, err = os.Stat(credsFile); os.IsNotExist(err) {
		log.Fatalf("%s: no such file or directory", credsFile)
	} else {
		f, _ := os.Open(credsFile)
		defer f.Close()

		byteValue, _ := ioutil.ReadAll(f)

		if err = json.Unmarshal(byteValue, &j); err != nil {
			log.Fatalf("%s: %v", credsFile, err)
		} else {
			conf.ClientID = j.Creds.ClientID
			conf.ClientSecret = j.Creds.ClientSecret
			conf.Endpoint.AuthURL = j.Creds.AuthURL
			conf.Endpoint.TokenURL = j.Creds.TokenURL
		}
	}
}

func saveToken(tokenFile string, token *oauth2.Token) {
	// can't create a folder
	f, err := os.OpenFile(tokenFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

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
