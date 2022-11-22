package auth

import (
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"
)

func getRandomPort() (string, error) {
	listener, err := net.Listen("tcp", "localhost:0")
	defer listener.Close()
	return ":" + strconv.Itoa(listener.Addr().(*net.TCPAddr).Port), err
}

func (resp *response) fetchFromHTTP() {
	listener, _ := net.Listen("tcp", config.RedirectURL[len("http://"):])
	svr := http.Server{}
	defer svr.Close()

	receivedFlag := make(chan bool)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if len(q.Get("error")) == 0 {
			io.WriteString(w, "qGmail has successfully received user's consent. You may close this tab.")
			resp.Code = q.Get("code")
			resp.State = q.Get("state")
			receivedFlag <- true
		} else {
			io.WriteString(w, "The user declined authorization request for qGmail.")
			receivedFlag <- false
		}
	})

	go func() {
		if err := svr.Serve(listener); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// Wait for an HTTP request for some time.
	select {
	case rf := <-receivedFlag:
		if !rf {
			log.Fatal("qGmail was not authorized by the user.")
		}
		// TODO: Make timeout configurable through a flag?
	case <-time.After(time.Duration(1) * time.Minute):
		log.Fatal("Error: Failed to authorize qGmail due to timeout.")
	}
}
