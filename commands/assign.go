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
	message := fmt.Sprintf("やったね <@%s> ちゃん :tada: %s のレビュアーになったまる", assigner(), extructTarget(ev.Text))
	if _, _, err := client.PostMessage(ev.Channel, message, params); err != nil {
		log.Printf("failed to post message: %s", err)
	}

}

func members() []string {
	return []string{"eccyan", "ai", "yuji.ueda", "ueshima", "yabu", "tkatagiri", "iizuka.daisuke"}
}

func normal() float64 {
	// Avg. assigned MR in a week, Std-dev ± 2
	return rand.NormFloat64()*2.0 + 4.0
}

func assigner() string {
	mem := members()
	return mem[int(normal()*1234567890)%len(mem)]
}

func extructTarget(text string) string {
	re := regexp.MustCompile(`assign\s+(.+)$`)
	match := re.FindStringSubmatch(text)
	if len(match) == 0 {
		return ""
	}

	return match[1]
}
