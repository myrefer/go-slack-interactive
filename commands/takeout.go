package commands

import (
	"fmt"
	api "github.com/nlopes/slack"
	"log"
	"time"

	"github.com/kelseyhightower/envconfig"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

type envConfig struct {
	Credentials   string `envconfig:"EQUIPMENT_MANAGEMENT_CREDENTIALS" required:"true"`
	SpreadSheetID string `envconfig:"EQUIPMENT_MANAGEMENT_SPREAD_SHEET_ID" required:"true"`
}

func Takeout(ev *api.MessageEvent, client *api.Client) {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Fatalf("[ERROR] Failed to process env var: %s", err)
	}

	config, err := google.JWTConfigFromJSON(
		[]byte(env.Credentials),
		"https://www.googleapis.com/auth/spreadsheets",
	)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	srv, err := sheets.New(config.Client(oauth2.NoContext))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	user, err := client.GetUserInfo(ev.User)
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
		return
	}

	var vr sheets.ValueRange
	vr.Values = append(vr.Values, []interface{}{time.Now(), user.Profile.RealName, user.Profile.Email, ev.Text})

	_, err = srv.Spreadsheets.Values.Append(env.SpreadSheetID, "A1", &vr).
		ValueInputOption("RAW").
		InsertDataOption("INSERT_ROWS").
		Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	params := api.NewPostMessageParameters()
	params.LinkNames = 1
	params.EscapeText = false
	message := fmt.Sprintf("<@%s> 了解まる :kira: 大切に使ってね！", ev.User)
	if _, _, err := client.PostMessage(ev.Channel, message, params); err != nil {
		log.Printf("failed to post message: %s", err)
	}
}
