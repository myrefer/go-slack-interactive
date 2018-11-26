package commands

import (
	"fmt"
	"log"

	api "github.com/nlopes/slack"
)

func Ping(ev *api.MessageEvent, client *api.Client) {
	log.Printf("ping")

	params := api.NewPostMessageParameters()
	params.LinkNames = 1
	params.EscapeText = false
	message := fmt.Sprintf("<@%s> pong", ev.User)
	if _, _, err := client.PostMessage(ev.Channel, message, params); err != nil {
		log.Printf("failed to post message: %s", err)
	}
}
