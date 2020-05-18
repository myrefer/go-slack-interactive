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
	IOS      = "ios"
	ANDROID  = "android"
	FRONT    = "front"
	BACK     = "back"
	ALL      = "all"
	RELEASE  = "release"
	TECHPERM = "tech-perm"
)

const (
	OK = 0
	NG = 1
)

const (
	SUMIYOSHI = "U5NHXB5T7"
	UESHIMA   = "U073ZTJS3"
	YABUSHITA = "UBHSQJ82E"
	NAKAYAMA  = "UF790B4H5"
	KATAGIRI  = "U06G2DMP1"
	KATO      = "UE75GQMFB"
	KAWANO    = "UF7J5SKT5"
	TOKUNAGA  = "USC6UBULE"
	MORI      = "UUQAEJABA"
	NAGAHARA  = "UUCS31K9R"
	AOKI      = "UUE7DFJ65"
	YANBE     = "U8WA4JY7N"
	SATO      = "UNYKNHP5M"
)

func Assign(ev *api.MessageEvent, client *api.Client) {
	rand.Seed(time.Now().Unix())

	params := api.NewPostMessageParameters()
	params.LinkNames = 1
	params.EscapeText = false

	var message string

	text, err := parseParams(ev.Text)
	if err != OK {
		message = fmt.Sprintf("< assign [ android | ios | front | back | release | tech-perm | all ] レビュー対象 > のフォーマットで話しかけてほしいまる")
	} else {
		assigner, err := assigner(text[2], ev.User)
		if err != OK {
			message = fmt.Sprintf("< assign [ android | ios | front | back | release | tech-perm | all ] レビュー対象 > のフォーマットで話しかけてほしいまる")
		} else {
			target := text[3]
			message = fmt.Sprintf("やったね <@%s> ちゃん :tada: %s のレビュアーになったまる!", assigner, target)
		}
	}
	if _, _, err := client.PostMessage(ev.Channel, message, params); err != nil {
		log.Printf("failed to post message: %s", err)
	}

}

// test
func members(category string) ([]string, int) {
	switch category {
	case IOS:
		return []string{SUMIYOSHI, KATAGIRI, KAWANO, MORI, YABUSHITA, NAKAYAMA}, OK
	case ANDROID:
		return []string{SUMIYOSHI, KATAGIRI, KAWANO, YABUSHITA, NAKAYAMA}, OK
	case FRONT:
		return []string{UESHIMA, KATO, TOKUNAGA}, OK
	case BACK:
		return []string{SUMIYOSHI, YABUSHITA, NAKAYAMA, KATAGIRI, NAGAHARA, AOKI}, OK
	case TECHPERM:
		return []string{SUMIYOSHI, UESHIMA, YABUSHITA, NAKAYAMA, TOKUNAGA}, OK
	case RELEASE:
		return []string{SUMIYOSHI, UESHIMA, YABUSHITA, NAKAYAMA, TOKUNAGA}, OK
	case ALL:
		return []string{SUMIYOSHI, UESHIMA, YABUSHITA, NAKAYAMA, KATAGIRI, KATO, KAWANO, TOKUNAGA, NAGAHARA, AOKI, MORI}, OK
	default:
		return nil, NG
	}
}

func normal() float64 {
	// Avg. assigned MR in a week, Std-dev ± 2
	return rand.NormFloat64()*2.0 + 4.0
}

func assigner(category string, user string) (string, int) {
	mem, err := members(category)
	if err != OK {
		return "", err
	}
	reviewer := mem[int(normal()*1234567890)%len(mem)]
	if reviewer == user {
		reviewer, err = assigner(category, user)
	}
	return reviewer, err
}

func parseParams(text string) ([]string, int) {
	match := strings.Fields(text)
	if len(match) < 3 {
		return nil, NG
	}
	return match, OK
}
