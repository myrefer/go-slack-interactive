package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/myrefer/go-slack-interactive/slack"
	api "github.com/nlopes/slack"
)

const (
	// action is used for slack attament action.
	actionSelect = "select"
	actionStart  = "start"
	actionCancel = "cancel"
)

type SlackListener struct {
	client    *api.Client
	botID     string
	channelID string
}

// LstenAndResponse listens slack events and response
// particular messages. It replies by slack message button.
func (s *SlackListener) ListenAndResponse() {
	rtm := s.client.NewRTM()

	// Start listening slack events
	go rtm.ManageConnection()

	// Handle slack events
	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *api.MessageEvent:
			if err := s.handleMessageEvent(ev); err != nil {
				log.Printf("[ERROR] Failed to handle message: %s", err)
			}
		}
	}
}

// handleMesageEvent handles message events.
func (s *SlackListener) handleMessageEvent(ev *api.MessageEvent) error {
	// Only response in specific channel. Ignore else.
	if ev.Channel != s.channelID {
		log.Printf("%s %s", ev.Channel, ev.Msg.Text)
		return nil
	}

	// Only response mention to bot. Ignore else.
	if !strings.HasPrefix(ev.Msg.Text, fmt.Sprintf("<@%s> ", s.botID)) {
		return nil
	}

	// Parse message
	mux := slack.NewServeMux()
	mux.Handle("hey", slack.HandlerFunc(Hey))
	mux.ServeMessage(ev, s.client)

	return nil
}

func Hey(ev *api.MessageEvent, client *api.Client) {
	log.Printf("Start Hey")
	// value is passed to message handler when request is approved.
	attachment := api.Attachment{
		Text:       "Which beer do you want? :beer:",
		Color:      "#f9a41b",
		CallbackID: "beer",
		Actions: []api.AttachmentAction{
			{
				Name: actionSelect,
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
				Name:  actionCancel,
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
	log.Printf("Succeeded to post")
}
