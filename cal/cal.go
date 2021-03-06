package cal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

type Server struct {
	Name        string
	CalendarURL string
}

type IntroEvent struct {
	Place string
	Event *calendar.Event
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func GetEvents() {
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := calendar.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}

	t := time.Now().Format(time.RFC3339)

	servers := []Server{
		Server{"IKM", "8skmlerhojthfhm8c8h76deqds@group.calendar.google.com"}, // IKM
		// "l9k5hhna9anjlrjmrum3v4m4hg@group.calendar.google.com", // Flash
		Server{"AIQ", "10fmjofa984qpol7oubsfs6eq8@group.calendar.google.com"},      // AIQ
		Server{"Tamayura", "mcr5e8d0dbe09s5u7orrqlg0co@group.calendar.google.com"}, // Tamayura
	}

	var introEvents []IntroEvent
	// var evelist [](*calendar.Event)
	for _, server := range servers {
		url := server.CalendarURL
		events, err := srv.Events.List(url).ShowDeleted(false).
			SingleEvents(true).TimeMin(t).MaxResults(20).OrderBy("startTime").Do()
		if err != nil {
			log.Fatalf("Unable to retrieve next ten of the user's events: %v", err)
		}
		for _, item := range events.Items {
			// evelist = append(evelist, item)
			introEvents = append(introEvents, IntroEvent{server.Name, item})
			// date := item.Start.DateTime
			// if date == "" {
			// 	date = item.Start.Date
			// }
			// fmt.Printf("%v (%v)\n", item.Summary, date)
		}
	}

	sort.Slice(introEvents, func(i, j int) bool { return compareEvent(introEvents[i].Event.Start, introEvents[j].Event.Start) })
	for _, ieve := range introEvents {
		eve := ieve.Event
		date := eve.Start.DateTime
		if date == "" {
			date = eve.Start.Date
		}
		fmt.Printf("%v - [%v]%v\n", date, ieve.Place, eve.Summary)
	}
}

func compareEvent(dt1, dt2 *calendar.EventDateTime) bool {
	dt1s := dt1.Date
	dt2s := dt2.Date
	if dt1s == "" {
		dt1s = dt1.DateTime
	}
	if dt2s == "" {
		dt2s = dt2.DateTime
	}
	return dt1s < dt2s
}
