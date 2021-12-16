package schemas

type IntroSchedule struct {
	No          int      `json:"no"`
	CircleId    string   `json:"circle"`
	EventId     string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Start       DateData `json:"start"`
	End         DateData `json:"end"`
	IsOffline   bool     `json:"offline"`
}

type DateData struct {
	Date string `json:"date"`
	Time string `json:"time"`
}

type Config struct {
	CheckDistance int          `json:"check_distance"`
	Msg           Message      `json:"message"`
	Dynamo        DynamoConfig `json:"dynamo"`
}

type Message struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type DynamoConfig struct {
	AccessToken string `json:"access_token"`
	Secret      string `json:"secret"`
}
