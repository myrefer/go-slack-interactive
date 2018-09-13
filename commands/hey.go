package commands

import (
	api "github.com/nlopes/slack"
	"log"
)

// TODO: Can I use namespace ? like below:
// HeyAction.Select
const (
	HeyCallbackID   = "beer"
	HeyActionSelect = "select"
	HeyActionStart  = "start"
	HeyActionCancel = "cancel"
)

func Hey(ev *api.MessageEvent, client *api.Client) {
	attachment := api.Attachment{
		Text:       "Which beer do you want? :beer:",
		Color:      "#f9a41b",
		CallbackID: HeyCallbackID,
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
