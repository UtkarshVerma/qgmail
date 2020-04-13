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
	"strconv"

	"golang.org/x/oauth2"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type (
	authParams struct {
		CodeVerifier string
		Request      authRequest
		Response     authResponse
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

func newOauthConf(credsFile string, conf *config) *oauth2.Config {
	oauthConf := &oauth2.Config{
		RedirectURL: "http://localhost:" + strconv.Itoa(conf.RedirectPort),
		Scopes:      auth.Request.Scopes,
	}
	readCredentials(credsFile, oauthConf)
	return oauthConf
}

func readCredentials(credsFile string, conf *oauth2.Config) {
	var err error
	var j struct {
		Creds *clientCreds `json:"installed"`
	}

	if _, err = os.Stat(credsFile); os.IsNotExist(err) {
		log.Fatalf("%s: no such file or directory", credsFile)
	} else {
		f, err := os.Open(credsFile)
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

func newAuthParams() *authParams {
	a := &authParams{
		CodeVerifier: generateRandomString(43),
		Request: authRequest{
			State:               generateRandomString(10),
			CodeChallengeMethod: "S256",
			Scopes:              []string{gmail.GmailLabelsScope},
		},
		Response: authResponse{},
	}
	a.Request.CodeChallenge = createCodeChallenge(a.CodeVerifier,
		a.Request.CodeChallengeMethod)
	return a
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	opts := []oauth2.AuthCodeOption{
		oauth2.SetAuthURLParam("access_type", "offline"),
		oauth2.SetAuthURLParam("code_challenge",
			auth.Request.CodeChallenge),
		oauth2.SetAuthURLParam("code_challenge_method",
			auth.Request.CodeChallengeMethod),
	}
	openURL(config.AuthCodeURL(auth.Request.State, opts...))
	// Fetch the authorization code.
	fmt.Println("Paste the authorization code:")
	if _, err := fmt.Scan(&auth.Response.Code); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	// Verify state
	// Handle user decline

	token, err := config.Exchange(context.TODO(), auth.Response.Code,
		oauth2.SetAuthURLParam("code_verifier", auth.CodeVerifier))
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

func saveToken(tokenFile string, token *oauth2.Token) {
	// can't create a folder
	f, err := os.OpenFile(tokenFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	defer f.Close()
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	json.NewEncoder(f).Encode(token)
}

// Retrieve token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	token := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(token)
	return token, err
}

func newGmailService(config *oauth2.Config, token *oauth2.Token) (*gmail.Service, error) {
	tokenSource := config.TokenSource(context.TODO(), token)
	service, err := gmail.NewService(context.TODO(), option.WithTokenSource(tokenSource))
	if err != nil {
		log.Fatalf("Unable to create Gmail service: %v", err)
	}
	return service, err
}
