package http

import (
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/utkarshverma/qgmail/config"
)

var listener net.Listener

func RandomPort() (string, error) {
	listener, err := net.Listen("tcp", "localhost:0")
	defer listener.Close()
	return ":" + strconv.Itoa(listener.Addr().(*net.TCPAddr).Port), err
}

func StartServer(address string, code, state *string) {
	listener, _ = net.Listen("tcp", address)
	svr := http.Server{}
	defer svr.Close()

	receivedFlag := make(chan bool)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()

		// Check if user has authorized qGmail.
		if len(q.Get("error")) == 0 {
			io.WriteString(w, "qGmail has successfully received user's consent. You may close this tab.")
			*code = q.Get("code")
			*state = q.Get("state")
			receivedFlag <- true
		} else {
			io.WriteString(w, "Error: The user declined the authorization request for qGmail.")
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
	case <-time.After(time.Duration(*config.Init.Timeout) * time.Minute):
		log.Fatal("Error: Failed to authorize qGmail due to timeout.")
	}
}
