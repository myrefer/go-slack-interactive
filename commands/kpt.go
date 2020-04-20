package commands

func Kpt(ev *api.MessageEvent, client *api.Client) {
	params := api.NewPostMessageParameters()
	params.LinkNames = 1
	params.EscapeText = false
	
	member := []string{SUMIYOSHI, UESHIMA, YABUSHITA, NAKAYAMA, KATAGIRI, KATO, KAWANO, CHIBA, TOKUNAGA, NAGAHARA, AOKI, MORI, YANBE, SATO}

	perm := []string{SUMIYOSHI, UESHIMA, YABUSHITA, NAKAYAMA, TOKUNAGA}

	var teamA []string
	var teamB []string
	var facilitatorA string
	var facilitatorB string
	var result bool

	for {
		shuffle(member)
		teamA = member[0 : len(member)/2]
		teamB = member[len(member)/2:]
		facilitatorA, result = assignFacilitatior(teamA, perm)
		if !result {
			continue
		}
		facilitatorB, result = assignFacilitatior(teamB, perm)
		if !result {
			continue
		}
		break
	}

	mentionA := generateMention(teamA)
	mentionB := generateMention(teamB)

	messageA := fmt.Sprintf("今日のKPTのチームAは %s まる！ ファシリテーターは <%s> まる！", mentionA, facilitatorA)
	messageB := fmt.Sprintf("今日のKPTのチームBは %s まる！ ファシリテーターは <%s> まる！", mentionB, facilitatorB)

	if _, _, err := client.PostMessage(ev.Channel, messageA, params); err != nil {
		log.Printf("failed to post message: %s", err)
	}
	if _, _, err := client.PostMessage(ev.Channel, messageB, params); err != nil {
		log.Printf("failed to post message: %s", err)
	}
}

func shuffle(data []string) {
	rand.Seed(time.Now().UnixNano())
	for i := range data {
		j := rand.Intn(i + 1)
		data[i], data[j] = data[j], data[i]
	}
}

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

func choice(s []string) string {
	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(len(s))
	return s[i]
}

func generateMention(team []string) string {
	var list []string
	for _, mem := range team {
		list = append(list, "<@"+mem+">")
	}
	message := strings.Join(list[:], " ")
	return message
}