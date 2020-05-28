package commands

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	api "github.com/nlopes/slack"
)

func Release(ev *api.MessageEvent, client *api.Client) {

	params := api.NewPostMessageParameters()
	params.LinkNames = 1
	params.EscapeText = false

	perm := []string{SUMIYOSHI, UESHIMA, YABUSHITA, NAKAYAMA, TOKUNAGA}

	web, ios, android := assign(perm)

	message := fmt.Sprintf("WEB担当：<@%s>\niOS担当：<@%s>\nAndroid担当：<@%s>\nヨロシクまる！", web, ios, android)

	if _, _, err := client.PostMessage(ev.Channel, message, params); err != nil {
		log.Printf("failed to post message: %s", err)
	}
}

func choice(s []string) string {
	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(len(s))
	return s[i]
}

func assign(s []string) (string, string, string) {
	web := choice(s)
	ios := choice(s)
	android := choise(s)
	if web == ios || web == android || ios == android {
		web, ios, android = assign(perm)
	}
	return web, ios, android
}
