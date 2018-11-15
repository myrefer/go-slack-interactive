package commands

import (
	"fmt"
	api "github.com/nlopes/slack"
	"log"
	"math/rand"
	"regexp"
	"time"
)

func Assign(ev *api.MessageEvent, client *api.Client) {
	rand.Seed(time.Now().Unix())

	params := api.NewPostMessageParameters()
	params.LinkNames = 1
	params.EscapeText = false
	message := fmt.Sprintf("おめでとうございます :tada: <@%s> が https://www.pivotaltracker.com/story/show/%s のレビュアーです。", assigner(), extructPID(ev.Text))
	if _, _, err := client.PostMessage(ev.Channel, message, params); err != nil {
		log.Printf("failed to post message: %s", err)
	}

}

func members() []string {
	return []string{"eccyan", "ai", "yuji.ueda", "ueshima", "yabu", "tkatagiri", "iizuka.daisuke"}
}

func assigner() string {
	mem := members()
	return mem[rand.Intn(len(mem))]
}

func extructPID(text string) string {
	re := regexp.MustCompile(`\s+#(\d+)`)
	match := re.FindStringSubmatch(text)
	if len(match) == 0 {
		return ""
	}

	return match[1]
}
