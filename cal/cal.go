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

type Circle struct {
	SimpleName  string         `json:"simple"`
	Name        string         `json:"name"`
	Overview    []string       `json:"overview"`
	Detail      []CircleDetail `json:"detail"`
	CalendarURL string         `json:"url"`
}

type CircleDetail struct {
	Item  string `json:"item"`
	Link  string `json:"link"`
	Value string `json:"value"`
}

type IntroEvent struct {
	CircleId string
	Event    *calendar.Event
}

type IntroSchedule struct {
	No          int      `json:"no"`
	CircleId    string   `json:"circle"`
	EventId     string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Start       DateData `json:"start"`
	End         DateData `json:"end"`
}

type DateData struct {
	Date string `json:"date"`
	Time string `json:"time"`
}

type ScheduleConfig struct {
	Update    string                   `json:"update"`
	Schedules map[string]IntroSchedule `json:"schedules"`
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

func getCircleFromJSON() (servers map[string]Circle) {
	bytes, err := ioutil.ReadFile("circles.json")
	if err != nil {
		log.Fatalf("Can't read `servers.json`. %v\n", err)
	}
	if err = json.Unmarshal(bytes, &servers); err != nil {
		log.Fatalf("Can't parse `servers.json`. %v\n", err)
	}
	return
}

func GetEvents() []IntroEvent {
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

	now := time.Now()
	t := now.Format(time.RFC3339)
	month3 := now.Add(time.Duration(90*24) * time.Hour)
	m3 := month3.Format(time.RFC3339)

	servers := getCircleFromJSON()

	var introEvents []IntroEvent
	// var evelist [](*calendar.Event)
	for key, server := range servers {
		url := server.CalendarURL
		events, err := srv.Events.List(url).ShowDeleted(false).
			SingleEvents(true).TimeMin(t).TimeMax(m3).OrderBy("startTime").Do()
		if err != nil {
			log.Fatalf("Unable to retrieve next ten of the user's events: %v\n", err)
		}
		for _, item := range events.Items {
			// evelist = append(evelist, item)
			introEvents = append(introEvents, IntroEvent{key, item})
			// date := item.Start.DateTime
			// if date == "" {
			// 	date = item.Start.Date
			// }
			// fmt.Printf("%v (%v)\n", item.Summary, date)
		}
	}

	sort.Slice(introEvents, func(i, j int) bool { return compareEvent(introEvents[i].Event.Start, introEvents[j].Event.Start) })
	return introEvents
	// for _, ieve := range introEvents {
	// 	eve := ieve.Event
	// 	date := eve.Start.DateTime
	// 	if date == "" {
	// 		date = eve.Start.Date
	// 	}
	// 	fmt.Printf("%v - [%v]%v\n", date, ieve.Place, eve.Summary)
	// }
}

func GetScheduleJson() ScheduleConfig {
	js, err := ioutil.ReadFile("schedule.json")
	if err != nil {
		log.Fatalf("Can't read schedule.json: %v\n", err)
	}
	var cfg ScheduleConfig
	if err := json.Unmarshal(js, &cfg); err != nil {
		log.Fatalf("Unmarshal error when parse schedule.json: %v\n", err)
	}
	return cfg
}

func MakeScheduleJson() {
	introEvents := GetEvents()
	schedules := map[string]IntroSchedule{}
	var prev string
	for idx, eve := range introEvents {
		var stime, sdate, etime, edate string
		if eve.Event.Start.Date == "" {
			start, err := time.Parse(time.RFC3339, eve.Event.Start.DateTime)
			if err != nil {
				log.Fatalf("Start time can't parse!: %v\n", err)
			}
			sdate = start.Format("2006/01/02(Mon)")
			stime = start.Format("15:04")
		} else {
			tmpsdate, err := time.Parse("2006-01-02", eve.Event.Start.Date)
			if err != nil {
				log.Fatalf("Start date can't parse!: %v\n", err)
			}
			sdate = tmpsdate.Format("2006/01/02(Mon)")
		}

		if eve.Event.End.Date == "" {
			end, err := time.Parse(time.RFC3339, eve.Event.End.DateTime)
			if err != nil {
				log.Fatalf("End time can't parse!: %v\n", err)
			}
			etime = end.Format("15:04")
		}
		if prev == sdate {
			sdate = ""
		} else {
			prev = sdate
		}
		schedules[eve.Event.Id] = IntroSchedule{
			No:          idx + 1,
			CircleId:    eve.CircleId,
			EventId:     eve.Event.Id,
			Title:       eve.Event.Summary,
			Description: eve.Event.Description,
			Start:       DateData{Date: sdate, Time: stime},
			End:         DateData{Date: edate, Time: etime},
		}
	}

	js, err := json.MarshalIndent(ScheduleConfig{Update: time.Now().Format("2006-01-02 15:04:05"), Schedules: schedules}, "", "    ")
	if err != nil {
		log.Fatalf("Json marshal error: %v\n", err)
	}
	ioutil.WriteFile("schedule.json", js, 0644)
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
