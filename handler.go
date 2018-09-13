package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/myrefer/go-slack-interactive/commands"
	api "github.com/nlopes/slack"
)

// interactionHandler handles interactive callback response.
type interactionHandler struct {
	slackClient       *api.Client
	verificationToken string
}

func (h interactionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Printf("[ERROR] Invalid method: %s", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("[ERROR] Failed to read request body: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonStr, err := url.QueryUnescape(string(buf)[8:])
	if err != nil {
		log.Printf("[ERROR] Failed to unespace request body: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var callback api.AttachmentActionCallback
	if err := json.Unmarshal([]byte(jsonStr), &callback); err != nil {
		log.Printf("[ERROR] Failed to decode json callback from slack: %s", jsonStr)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Only accept callback from slack with valid token
	if callback.Token != h.verificationToken {
		log.Printf("[ERROR] Invalid token: %s", callback.Token)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if callback.CallbackID != commands.HeyCallbackID {
		log.Printf("[ERROR] Invalid callbackId: %s", callback.CallbackID)
		return
	}

	hey := commands.NewHey(callback.CallbackID)
	hey.ServeInteractiveAction(&callback, w)
}
