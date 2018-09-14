package commands

import (
	"fmt"
	api "github.com/nlopes/slack"
	"log"
)

func Ping(ev *api.MessageEvent, client *api.Client) {

	params := api.NewPostMessageParameters()
	params.LinkNames = 1
	params.EscapeText = false
	message := fmt.Sprintf("<@%s> pong", ev.User)
	if _, _, err := client.PostMessage(ev.Channel, message, params); err != nil {
		log.Printf("failed to post message: %s", err)
	}
}
