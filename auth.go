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
	oauthConf := &oauth2.Config{
		RedirectURL: "http://localhost:3000",
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
	openURL(config.AuthCodeURL(auth.Request.State, opts...))
	// will the localhost server ever fail?
	fmt.Println("Opening browser for user consent...")
	time.Sleep(2 * time.Second)
	auth.Response.getAuthCode("localhost:3000")

	// Verify state parameter
	if auth.Request.State != auth.Response.State {
		log.Fatal("Error: This request wasn't initialised by me.")
	}

	token, err := config.Exchange(context.TODO(), auth.Response.Code,
		oauth2.SetAuthURLParam("code_verifier", auth.CodeVerifier))
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return token
}

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
		fmt.Printf("Go to the following link in your browser then paste the "+
			"authorization code here: \n%v\n", url)
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
