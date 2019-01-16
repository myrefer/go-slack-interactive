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
	RELEASE   = "release"
	INTERVIEW = "interview"
)

const (
	OK = 0
	NG = 1
)

const (
	UEDA      = "U0A1HHDQR"
	SUMIYOSHI = "U5NHXB5T7"
	UESHIMA   = "U073ZTJS3"
	YABUSHITA = "UBHSQJ82E"
	NAKAYAMA  = "UF790B4H5"
	KATAGIRI  = "U06G2DMP1"
	IIZUKA    = "UC1H5MTRV"
	KATO      = "UE75GQMFB"
	KAWANO    = "UF7J5SKT5"
	MORI      = "UF7HXU4HH"
)

func Assign(ev *api.MessageEvent, client *api.Client) {
	rand.Seed(time.Now().Unix())

	params := api.NewPostMessageParameters()
	params.LinkNames = 1
	params.EscapeText = false
	var message string

	text, err := parseParams(ev.Text)
	if err != OK {
		message = fmt.Sprintf("< assign [ apl | front | back | release | interview | all ] レビュー対象 > のフォーマットで話しかけてほしいまる")
	} else {
		assigner, err := assigner(text[2])
		if err != OK {
			message = fmt.Sprintf("< assign [ apl | front | back | release | interview | all ] レビュー対象 > のフォーマットで話しかけてほしいまる")
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
		return []string{UEDA, SUMIYOSHI, KATAGIRI, IIZUKA, KAWANO}, OK
	case FRONT:
		return []string{UEDA, SUMIYOSHI, UESHIMA, KATAGIRI, IIZUKA, KATO, MORI}, OK
	case BACK:
		return []string{UEDA, SUMIYOSHI, YABUSHITA, NAKAYAMA, KATAGIRI, IIZUKA, MORI, KAWANO}, OK
	case INTERVIEW:
		return []string{UEDA, SUMIYOSHI, UESHIMA, YABUSHITA}, OK
	case RELEASE:
		return []string{UEDA, SUMIYOSHI, UESHIMA, YABUSHITA, NAKAYAMA}, OK
	case ALL:
		return []string{UEDA, SUMIYOSHI, UESHIMA, YABUSHITA, NAKAYAMA, KATAGIRI, IIZUKA, KATO, MORI, KAWANO}, OK
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
