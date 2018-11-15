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
	text, channel := parse(ev.Text)
	message := fmt.Sprintf("%s", text)
	if _, _, err := client.PostMessage(channel, message, params); err != nil {
		log.Printf("failed to post message: %s", err)
	}

}

func parse(text string) (string, string) {
	re := regexp.MustCompile(`proxy\s+(.+)\s+#(.+)$`)
	match := re.FindStringSubmatch(text)
	if len(match) < 2 {
		return "", ""
	}
	if len(match) < 3 {
		return match[1], ""
	}

	return match[1], match[2]
}
