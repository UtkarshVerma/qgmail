package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"golang.org/x/oauth2"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"

	"github.com/UtkarshVerma/qgmail/pkce"
)

type (
	request struct {
		State               string
		CodeChallenge       string
		CodeChallengeMethod string
		Scopes              []string
	}

	response struct {
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
	// CredsFile stores path to Gmail API credentials.
	CredsFile string

	config = &oauth2.Config{}
	token  *oauth2.Token
	params = &struct {
		CodeVerifier string
		Request      request
		Response     response
	}{
		CodeVerifier: pkce.RandomString(43),
		Request: request{
			State:               pkce.RandomString(10),
			CodeChallengeMethod: "S256",
			Scopes:              []string{gmail.GmailLabelsScope},
		},
	}
)

func init() {
	params.Request.CodeChallenge = pkce.CodeChallenge(params.CodeVerifier,
		params.Request.CodeChallengeMethod)

	redirectPort, err := getRandomPort()
	if err != nil {
		log.Fatal(err)
	}

	config = &oauth2.Config{
		RedirectURL: "http://localhost" + redirectPort,
		Scopes:      params.Request.Scopes,
	}
}

// NewGmailService creates a new `gmail.Service` using `Config` and `Token`.
func NewGmailService() (*gmail.Service, error) {
	tokenSource := config.TokenSource(context.TODO(), token)
	service, err := gmail.NewService(context.TODO(),
		option.WithTokenSource(tokenSource))
	return service, err
}

// ReadCredentials unmarshals `CredsFile` to `config`.
func ReadCredentials() error {
	var err error
	var j struct {
		Creds *clientCreds `json:"installed"`
	}
	if _, err = os.Stat(CredsFile); os.IsNotExist(err) {
		return err
	}

	f, _ := os.Open(CredsFile)
	defer f.Close()

	byteValue, _ := ioutil.ReadAll(f)

	if err = json.Unmarshal(byteValue, &j); err != nil {
		return err
	}

	config.ClientID = j.Creds.ClientID
	config.ClientSecret = j.Creds.ClientSecret
	config.Endpoint = oauth2.Endpoint{
		AuthURL:   j.Creds.AuthURL,
		TokenURL:  j.Creds.TokenURL,
		AuthStyle: oauth2.AuthStyleAutoDetect,
	}
	return nil
}

// GetToken requests an authorization token and stores it in `Token`.
func GetToken() (err error) {
	fetchAuthCode()

	// Exchange token with authorization code.
	token, err = config.Exchange(context.TODO(), params.Response.Code,
		oauth2.SetAuthURLParam("code_verifier", params.CodeVerifier))
	return err
}

func fetchAuthCode() {
	opts := []oauth2.AuthCodeOption{
		oauth2.SetAuthURLParam("access_type", "offline"),
		oauth2.SetAuthURLParam("code_challenge",
			params.Request.CodeChallenge),
		oauth2.SetAuthURLParam("code_challenge_method",
			params.Request.CodeChallengeMethod),
	}
	authURL := config.AuthCodeURL(params.Request.State, opts...)

	fmt.Printf("Open the following link in your web browser:\n%s\n\n",
		authURL)
	params.Response.fetchFromHTTP()
}

// SaveToken saves `Token` at `tokenPath`.
func SaveToken(tokenPath string) error {
	tokenFile, _ := json.Marshal(*token)

	// Create parent folders if non-existent.
	tokenDir := path.Dir(tokenPath)
	if _, err := os.Stat(tokenDir); os.IsNotExist(err) {
		os.MkdirAll(tokenDir, os.ModePerm)
	}

	return ioutil.WriteFile(tokenPath, tokenFile, 0600)
}

// ReadToken reads `tokenFile` to `Token`.
func ReadToken(tokenFile string) error {
	f, err := os.Open(tokenFile)
	if err != nil {
		return err
	}
	defer f.Close()
	if token == nil {
		token = &oauth2.Token{}
	}
	err = json.NewDecoder(f).Decode(token)
	return err
}
