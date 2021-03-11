package commands

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	api "github.com/nlopes/slack"
)

func Kpt2(ev *api.MessageEvent, client *api.Client) {

	params := api.NewPostMessageParameters()
	params.LinkNames = 1
	params.EscapeText = false

	//member := []string{SUMIYOSHI, UESHIMA, NAKAYAMA, KATAGIRI, MORI, YANBE, KAMINAGA, TAKADA, TOUYAMA, IDA, KOSHIMIZU, KATO}
	//perm := []string{SUMIYOSHI, UESHIMA, NAKAYAMA, KAMINAGA, KOSHIMIZU}
	teamA := []string{YABUSHITA, UESHIMA, KATO, NAKAYAMA, TOUYAMA, YANBE}
	permA := []string{YABUSHITA, UESHIMA, NAKAYAMA}
	teamB := []string{SUMIYOSHI, KAMINAGA, KATAGIRI, KOSHIMIZU, TAKADA, MORI}
	permB := []string{SUMIYOSHI, KAMINAGA, KOSHIMIZU}

	facilitatorA := choice(permA)
	secretaryA := assignSecretary(teamA, facilitatorA)

	message := fmt.Sprintf("今日のKPTのTeamAファシリテーターは <@%s> まる！ 書記は <@%s> まる！", facilitatorA, secretaryA)
	if _, _, err := client.PostMessage(ev.Channel, message, params); err != nil {
		log.Printf("failed to post message: %s", err)
	}

	facilitatorB := choice(permB)
	secretaryB := assignSecretary(teamB, facilitatorB)

	message = fmt.Sprintf("今日のKPTのTeamAファシリテーターは <@%s> まる！ 書記は <@%s> まる！", facilitatorB, secretaryB)
	if _, _, err := client.PostMessage(ev.Channel, message, params); err != nil {
		log.Printf("failed to post message: %s", err)
	}
}

func choice(s []string) string {
	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(len(s))
	return s[i]
}

func assignSecretary(s []string, facilitator string) string {
	secretary := choice(s)
	if secretary == facilitator {
		secretary = assignSecretary(s)
	}
	return secretary
}
