package auth

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

	"github.com/utkarshverma/qgmail/config"
	"github.com/utkarshverma/qgmail/http"
	"github.com/utkarshverma/qgmail/pkce"
	"golang.org/x/oauth2"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type (
	params struct {
		MustPaste    *bool
		MustShowURL  *bool
		CodeVerifier string
		Request      request
		Response     response
	}

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
	Config = &oauth2.Config{}
	Token  *oauth2.Token
	creds  = &clientCreds{}
	Params = &params{
		MustPaste:   config.Init.MustPaste,
		MustShowURL: config.Init.MustShowURL,

		CodeVerifier: pkce.RandomString(43),
		Request: request{
			State:               pkce.RandomString(10),
			CodeChallengeMethod: "S256",
			Scopes:              []string{gmail.GmailLabelsScope},
		},
	}
)

func init() {
	Params.Request.CodeChallenge = pkce.CodeChallenge(Params.CodeVerifier,
		Params.Request.CodeChallengeMethod)
	Config = newConfig(Params)
	readCredentials(**config.CredsFile, Config)
}

func NewGmailService(conf *oauth2.Config, token *oauth2.Token) (*gmail.Service, error) {
	tokenSource := conf.TokenSource(context.TODO(), token)
	service, err := gmail.NewService(context.TODO(),
		option.WithTokenSource(tokenSource))
	if err != nil {
		log.Fatalf("Unable to create Gmail service: %v", err)
	}
	return service, err
}

func newConfig(params *params) *oauth2.Config {
	var redirectURL string
	mustPaste := params.MustPaste

	// Use manual copy/paste method if error occurs.
	redirectPort, err := http.RandomPort()
	if err != nil {
		*mustPaste = true
	}
	if *mustPaste {
		redirectURL = "urn:ietf:wg:oauth:2.0:oob"
	} else {
		redirectURL = "http://localhost" + redirectPort
	}
	conf := &oauth2.Config{
		RedirectURL: redirectURL,
		Scopes:      params.Request.Scopes,
	}
	return conf
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
			conf.Endpoint = oauth2.Endpoint{
				AuthURL:   j.Creds.AuthURL,
				TokenURL:  j.Creds.TokenURL,
				AuthStyle: oauth2.AuthStyleAutoDetect,
			}
		}
	}
}

func GetToken(conf *oauth2.Config, token **oauth2.Token, params *params) {
	mustPaste := params.MustPaste
	opts := []oauth2.AuthCodeOption{
		oauth2.SetAuthURLParam("access_type", "offline"),
		oauth2.SetAuthURLParam("code_challenge",
			params.Request.CodeChallenge),
		oauth2.SetAuthURLParam("code_challenge_method",
			params.Request.CodeChallengeMethod),
	}

	// Fetch authorization code.
	authURL := conf.AuthCodeURL(params.Request.State, opts...)
	getAuthCode(conf, authURL, params)

	// Verify state parameter.
	if !*mustPaste && (params.Request.State != params.Response.State) {
		log.Fatal("Error: This request wasn't initialised by qGmail.")
	}

	// Exchange token with authorization code.
	tok, err := conf.Exchange(context.TODO(), params.Response.Code,
		oauth2.SetAuthURLParam("code_verifier", params.CodeVerifier))
	if err != nil {
		log.Fatalf("Error: Unable to retrieve token from the web.\n%v", err)
	} else {
		*token = tok
		fmt.Printf("%+v", *Token)
		fmt.Println("qGmail has been successfully authorized.")
	}
}

func getAuthCode(conf *oauth2.Config, authURL string, params *params) {
	mustPaste, mustShowURL := params.MustPaste, params.MustShowURL
	code, state := &params.Response.Code, &params.Response.State
	if *mustShowURL {
		fmt.Printf("Open the following link in your web browser:\n%s\n\n",
			authURL)
	} else {
		fmt.Println("Opening browser for user consent...")
		openURL(authURL)
		time.Sleep(2 * time.Second)
	}
	if *mustPaste {
		fmt.Println("Paste the authorization code here:")
		if _, err := fmt.Scan(code); err != nil {
			fmt.Println()
			log.Fatalf("Unable to read authorization code: %v", err)
		}
		fmt.Println()
	} else {
		http.StartServer(conf.RedirectURL[7:], code, state)
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

func SaveToken(tokenFile string, token *oauth2.Token) {
	// can't create a folder
	f, err := os.OpenFile(tokenFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Error: Unable to cache oauth token.\n%v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func ReadToken(tokenFile string, token **oauth2.Token) error {
	f, err := os.Open(tokenFile)
	if err != nil {
		return err
	}
	defer f.Close()
	if *token == nil {
		*token = &oauth2.Token{}
	}
	err = json.NewDecoder(f).Decode(token)
	return err
}
