package commands

import (
	"fmt"
	api "github.com/nlopes/slack"
	"log"
	"math/rand"
	"regexp"
	"time"
)

func Proxy(ev *api.MessageEvent, client *api.Client) {
	rand.Seed(time.Now().Unix())

	params := api.NewPostMessageParameters()
	params.LinkNames = 1
	params.EscapeText = false
	text := parse(ev.Text)
	message := fmt.Sprintf("%s", text)
	if _, _, err := client.PostMessage(ev.Channel, message, params); err != nil {
		log.Printf("failed to post message: %s", err)
	}

}

func parse(text string) string {
	re := regexp.MustCompile(`proxy\s+(.+)$`)
	match := re.FindStringSubmatch(text)
	if len(match) < 2 {
		return ""
	}

	return match[1]
}
