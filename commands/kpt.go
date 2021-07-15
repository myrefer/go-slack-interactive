package commands

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	api "github.com/nlopes/slack"
)

func Kpt(ev *api.MessageEvent, client *api.Client) {

	params := api.NewPostMessageParameters()
	params.LinkNames = 1
	params.EscapeText = false

	member := []string{SUMIYOSHI, UESHIMA, NAKAYAMA, KATAGIRI, YANBE, KAMINAGA, TAKADA, TOUYAMA, KOSHIMIZU}

	perm := []string{SUMIYOSHI, UESHIMA, NAKAYAMA, KAMINAGA, KOSHIMIZU}

	facilitator := choice(perm)
	secretary := assignSecretary(member, facilitator)

	message := fmt.Sprintf("今日のKPTのファシリテーターは <@%s> まる！ 書記は <@%s> まる！", facilitator, secretary)

	if _, _, err := client.PostMessage(ev.Channel, message, params); err != nil {
		log.Printf("failed to post message: %s", err)
	}
}
func Kpt2(ev *api.MessageEvent, client *api.Client) {

	params := api.NewPostMessageParameters()
	params.LinkNames = 1
	params.EscapeText = false

	//member := []string{SUMIYOSHI, UESHIMA, NAKAYAMA, KATAGIRI, YANBE, KAMINAGA, TAKADA, TOUYAMA, IDA, KOSHIMIZU, KATO}
	//perm := []string{SUMIYOSHI, UESHIMA, NAKAYAMA, KAMINAGA, KOSHIMIZU}
	teamA := []string{KATO, NAKAYAMA, KISHIMOTO, UESHIMA, KOSHIMIZU, TOUYAMA, TAKADA, KOBAYASHI}
	permA := []string{NAKAYAMA, KISHIMOTO, KOSHIMIZU, UESHIMA}
	teamB := []string{SUMIYOSHI, KAMINAGA, KATAGIRI, YSATO, YANBE, MASAKA, MANGOKU}
	permB := []string{SUMIYOSHI, KAMINAGA, YSATO, MASAKA, MANGOKU}

	facilitatorA := choice(permA)
	secretaryA := assignSecretary(teamA, facilitatorA)

	message := fmt.Sprintf("今日のKPTのTeamAファシリテーターは <@%s> まる！ 書記は <@%s> まる！", facilitatorA, secretaryA)
	if _, _, err := client.PostMessage(ev.Channel, message, params); err != nil {
		log.Printf("failed to post message: %s", err)
	}

	facilitatorB := choice(permB)
	secretaryB := assignSecretary(teamB, facilitatorB)

	message = fmt.Sprintf("今日のKPTのTeamBファシリテーターは <@%s> まる！ 書記は <@%s> まる！", facilitatorB, secretaryB)
	if _, _, err := client.PostMessage(ev.Channel, message, params); err != nil {
		log.Printf("failed to post message: %s", err)
	}
}

/**
func shuffle(data []string) {
	rand.Seed(time.Now().UnixNano())
	for i := range data {
		j := rand.Intn(i + 1)
		data[i], data[j] = data[j], data[i]
	}
}
*/

/**
func assignFacilitatior(team []string, perms []string) (string, bool) {
	candidates := []string{}
	for _, member := range team {
		for _, perm := range perms {
			if member == perm {
				candidates = append(candidates, member)
			}
		}
	}
	if len(candidates) == 0 {
		return "", false
	}
	facilitator := choice(candidates)

	return facilitator, true
}
*/

func choice(s []string) string {
	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(len(s))
	return s[i]
}

func assignSecretary(s []string, facilitator string) string {
	secretary := choice(s)
	if secretary == facilitator {
		secretary = assignSecretary(s, facilitator)
	}
	return secretary
}
