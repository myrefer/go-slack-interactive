package commands

import (
	"encoding/json"
	"fmt"
	"github.com/myrefer/go-slack-interactive/slack"
	api "github.com/nlopes/slack"
	"log"
	"net/http"
	"strings"
)

// TODO: Can I use namespace ? like below:
// HeyAction.Select
const (
	HeyCallbackID   = "beer"
	HeyActionSelect = "select"
	HeyActionStart  = "start"
	HeyActionCancel = "cancel"
)

type Hey struct {
	callbackID string
	mux        slack.ServeInteractiveActionMux
}

func NewHey(callbackID string) *Hey {
	h := new(Hey)
	h.callbackID = callbackID

	mux := slack.NewServeInteractiveActionMux(callbackID)
	mux.Handle("select", slack.InteractiveActionHandlerFunc(actionSelect))
	mux.Handle("start", slack.InteractiveActionHandlerFunc(actionStart))
	mux.Handle("cancel", slack.InteractiveActionHandlerFunc(actionCancel))
	return h
}

func (hey *Hey) ServeMessage(ev *api.MessageEvent, client *api.Client) {
	attachment := api.Attachment{
		CallbackID: hey.callbackID,
		Text:       "Which beer do you want? :beer:",
		Color:      "#f9a41b",
		Actions: []api.AttachmentAction{
			{
				Name: HeyActionSelect,
				Type: "select",
				Options: []api.AttachmentActionOption{
					{
						Text:  "Asahi Super Dry",
						Value: "Asahi Super Dry",
					},
					{
						Text:  "Kirin Lager Beer",
						Value: "Kirin Lager Beer",
					},
					{
						Text:  "Sapporo Black Label",
						Value: "Sapporo Black Label",
					},
					{
						Text:  "Suntory Malts",
						Value: "Suntory Malts",
					},
					{
						Text:  "Yona Yona Ale",
						Value: "Yona Yona Ale",
					},
					{
						Text:  "Plemiun Molt",
						Value: "Plemiun Molt",
					},
					{
						Text:  "Yebisu",
						Value: "Yebisu",
					},
				},
			},

			{
				Name:  HeyActionCancel,
				Text:  "Cancel",
				Type:  "button",
				Style: "danger",
			},
		},
	}

	params := api.PostMessageParameters{
		Attachments: []api.Attachment{
			attachment,
		},
	}

	if _, _, err := client.PostMessage(ev.Channel, "", params); err != nil {
		log.Printf("failed to post message: %s", err)
	}
}

func (hey *Hey) ServeInteractiveAction(callback *api.AttachmentActionCallback, w http.ResponseWriter) {
	hey.mux.ServeInteractiveAction(callback, w)
}

func actionSelect(callback *api.AttachmentActionCallback, w http.ResponseWriter) {
	action := callback.Actions[0]
	value := action.SelectedOptions[0].Value

	// Overwrite original drop down message.
	originalMessage := callback.OriginalMessage
	originalMessage.Attachments[0].Text = fmt.Sprintf("OK to order %s ?", strings.Title(value))
	originalMessage.Attachments[0].Actions = []api.AttachmentAction{
		{
			Name:  HeyActionStart,
			Text:  "Yes",
			Type:  "button",
			Value: "start",
			Style: "primary",
		},
		{
			Name:  HeyActionCancel,
			Text:  "No",
			Type:  "button",
			Style: "danger",
		},
	}

	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&originalMessage)
}

func actionStart(callback *api.AttachmentActionCallback, w http.ResponseWriter) {
	title := ":ok: your order was submitted! yay!"
	responseMessage(w, callback.OriginalMessage, title, "")
}

func actionCancel(callback *api.AttachmentActionCallback, w http.ResponseWriter) {
	title := fmt.Sprintf(":x: @%s canceled the request", callback.User.Name)
	responseMessage(w, callback.OriginalMessage, title, "")
}

// responseMessage response to the original slackbutton enabled message.
// It removes button and replace it with message which indicate how bot will work
func responseMessage(w http.ResponseWriter, original api.Message, title, value string) {
	original.Attachments[0].Actions = []api.AttachmentAction{} // empty buttons
	original.Attachments[0].Fields = []api.AttachmentField{
		{
			Title: title,
			Value: value,
			Short: false,
		},
	}

	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&original)
}
