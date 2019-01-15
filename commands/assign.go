package commands

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	api "github.com/nlopes/slack"
)

const (
	APL       = "apl"
	FRONT     = "front"
	BACK      = "back"
	ALL       = "all"
	INTERVIEW = "interview"
)

const (
	OK = 0
	NG = 1
)

func Assign(ev *api.MessageEvent, client *api.Client) {
	rand.Seed(time.Now().Unix())

	params := api.NewPostMessageParameters()
	params.LinkNames = 1
	params.EscapeText = false
	var message string

	text, err := parseParams(ev.Text)
	if err != OK {
		message = fmt.Sprintf("< assign [ apl | front | back | interview | all ] レビュー対象 > のフォーマットで話しかけてほしいまる")
	} else {
		assigner, err := assigner(text[2])
		if err != OK {
			message = fmt.Sprintf("< assign [ apl | front | back | interview | all ] レビュー対象 > のフォーマットで話しかけてほしいまる")
		} else {
			target := text[3]
			message = fmt.Sprintf("やったね <@%s> ちゃん :tada: %s のレビュアーになったまる", assigner, target)
		}
	}
	if _, _, err := client.PostMessage(ev.Channel, message, params); err != nil {
		log.Printf("failed to post message: %s", err)
	}

}

// test
func members(category string) ([]string, int) {
	switch category {
	case APL:
		return []string{"ai", "yuji.ueda", "tkatagiri", "iizuka.daisuke", "Takafumi Kawano"}, OK
	case FRONT:
		return []string{"ai", "yuji.ueda", "tkatagiri", "iizuka.daisuke", "ueshima", "yabu", "hitoshi.kato", "中山優", "Norihiro Mori"}, OK
	case BACK:
		return []string{"ai", "yuji.ueda", "tkatagiri", "iizuka.daisuke", "yabu", "Norihiro Mori", "中山優", "Takafumi Kawano"}, OK
	case INTERVIEW:
		return []string{"ai", "yuji.ueda", "ueshima", "yabu"}, OK
	case ALL:
		return []string{"ai", "yuji.ueda", "tkatagiri", "iizuka.daisuke", "yabu", "ueshima", "hitoshi.kato", "中山優", "Norihiro Mori", "Takafumi Kawano"}, OK
	default:
		return nil, NG
	}
}

func normal() float64 {
	// Avg. assigned MR in a week, Std-dev ± 2
	return rand.NormFloat64()*2.0 + 4.0
}

func assigner(category string) (string, int) {
	mem, err := members(category)
	if err != OK {
		return "", err
	}
	return mem[int(normal()*1234567890)%len(mem)], err
}

func parseParams(text string) ([]string, int) {
	match := strings.Fields(text)
	if len(match) < 3 {
		return nil, NG
	}
	return match, OK
}
