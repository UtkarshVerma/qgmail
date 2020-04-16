package main

import (
	"io"
	"log"
	"net/http"
)

// Handle user decline
func (a *authResponse) getAuthCode(address string) {
	svr := http.Server{Addr: address}
	defer svr.Close()

	receivedFlag := make(chan bool)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "qGmail has successfully received user's consent. You may close this tab.")
		a.Code = r.URL.Query().Get("code")
		a.State = r.URL.Query().Get("state")
		receivedFlag <- true
	})

	go func() {
		if err := svr.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// Wait for an HTTP request.
	<-receivedFlag
}
