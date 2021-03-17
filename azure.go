package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/gin-gonic/gin"

	"./cal"
)

func main() {
	if len(os.Args) > 1 {
		if os.Args[1] == "--get" {
			cal.MakeScheduleJson()
			log.Println("Complete to get events from Google Calendar.")
			return
		}
	}

	log.Println("Start Introquiz Portal Square Azure")

	go loopGet()

	// routerの初期設定
	router := gin.Default()

	// js,css,faviconなどを読み込むためのasstes設定
	router.LoadHTMLGlob("view/*.html")
	router.Static("/resource", "./resource")
	router.StaticFile("/favicon.ico", "./resource/favicon.ico")

	router.GET("/", func(ctx *gin.Context) {
		cfg := cal.GetScheduleJson()
		ctx.HTML(http.StatusOK, "main.html", gin.H{
			"update":  cfg.Update,
			"sche":    getTemplateSche(cfg.Schedules),
			"circles": getCircleList(),
		})
	})

	router.Run(":56417")

	log.Println("End Introquiz Portal Square Azure")
}

func loopGet() {
	for {
		cal.MakeScheduleJson()
		log.Println("Complete to get events from Google Calendar.")

		time.Sleep(time.Minute * 5)
	}
}

type Circle struct {
	Simple   string
	Name     string
	Overview string
	Twitter  string
	Contact  string
	Url      string
}

type TmpSchedule struct {
	Schedule cal.IntroSchedule
	EventId  string
	Simple   string
}

func getCircles() map[string]Circle {
	js, err := ioutil.ReadFile("circles.json")
	if err != nil {
		log.Fatalf("Can't read circles.json: %v\n", err)
	}

	var circles map[string]Circle
	err = json.Unmarshal(js, &circles)
	if err != nil {
		log.Fatalf("Unmarshal error circles.json: %v\n", err)
	}
	return circles
}

func getTemplateSche(sche map[string]cal.IntroSchedule) (sc []TmpSchedule) {
	cir := getCircles()
	for key, s := range sche {
		sc = append(sc, TmpSchedule{
			Schedule: s,
			EventId:  key,
			Simple:   cir[s.CircleId].Simple,
		})
	}
	sort.Slice(sc, func(i, j int) bool { return sc[i].Schedule.No < sc[j].Schedule.No })
	return sc
}

func getCircleList() (clist []string) {
	circles := getCircles()
	for _, c := range circles {
		clist = append(clist, c.Name)
	}
	return
}
