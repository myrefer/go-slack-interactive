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

	member := []string{SUMIYOSHI, UESHIMA, NAKAYAMA, KATAGIRI, KATO, TOKUNAGA, MORI, YANBE, SATO, KAMINAGA, TAKADA}

	perm := []string{SUMIYOSHI, UESHIMA, NAKAYAMA, TOKUNAGA, KAMINAGA}

	facilitator := choice(perm)
	secretary := assignSecretary(member, facilitator)

	message := fmt.Sprintf("今日のKPTのファシリテーターは <@%s> まる！ 書記は <@%s> まる！", facilitator, secretary)

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
		secretary = choice(s)
	}
	return secretary
}

/**
func generateMention(team []string) string {
	var list []string
	for _, mem := range team {
		list = append(list, "<@"+mem+">")
	}
	message := strings.Join(list[:], " ")
	return message
}
*/
