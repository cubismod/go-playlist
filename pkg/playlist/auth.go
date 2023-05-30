// Taken from https://github.com/zmb3/spotify/blob/master/examples/authenticate/authcode/authenticate.go
// ZMB3 is under an Apache License
package playlist

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/apex/log"
	"github.com/google/uuid"
	spotifyauth "github.com/zmb3/spotify/v2/auth"

	"github.com/zmb3/spotify/v2"
)

var redirectURI = generateUrl()

var (
	auth  = spotifyauth.New(spotifyauth.WithRedirectURL(redirectURI), spotifyauth.WithScopes(spotifyauth.ScopePlaylistModifyPrivate, spotifyauth.ScopePlaylistModifyPublic))
	ch    = make(chan *spotify.Client)
	state = uuid.New().String()
)

func generateUrl() string {
	host := "http://localhost"
	port := "8080"

	if os.Getenv("HOSTNAME") != "" {
		host = os.Getenv("HOSTNAME")
	}

	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	return fmt.Sprintf("%s:%s/callback", host, port)
}

func RunAuthServer() (*spotify.Client, error) {
	// first start an HTTP server
	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {})
	go func() {
		err := http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), nil)
		if err != nil {
			log.Fatalf("Unable to start web server", err)
		}
	}()

	url := auth.AuthURL(state)
	msg := fmt.Sprintf("Please log in to Spotify by visiting the following page in your browser: %s", url)

	req, _ := http.NewRequest("POST", os.Getenv("NTFY_URL"), strings.NewReader(msg))
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", os.Getenv("NTFY_PW")))
	req.Header.Set("Actions", fmt.Sprintf("view, Open, %s", url))
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		log.WithError(err).Error("Unable to post to Gotify server")
	}

	log.WithFields(log.Fields{
		"code": resp.StatusCode,
		"body": resp.Status,
	}).Info("Gotify req")

	log.Info(msg)

	// wait for auth to complete
	client := <-ch

	// use the client to make calls that require authorization
	user, err := client.CurrentUser(context.Background())
	if err != nil {
		log.WithError(err).Error("Unable to login")
	}
	fmt.Println("You are logged in as:", user.ID)

	return client, err
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(r.Context(), state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Error(err.Error())
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Errorf("State mismatch: %s != %s\n", st, state)
	}

	// use the token to get an authenticated client
	client := spotify.New(auth.Client(r.Context(), tok))
	fmt.Fprintf(w, "Login Completed!")
	ch <- client
}
