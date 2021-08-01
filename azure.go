package main

import (
	"encoding/json"
	"html/template"
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
			callSchedule()
			return
		}
	}

	log.Println("Start Introquiz Portal Square Azure")

	go loopGet()

	// routerの初期設定
	router := gin.Default()

	router.SetFuncMap(template.FuncMap{
		"safe": noescape,
	})
	// js,css,faviconなどを読み込むためのasstes設定
	router.LoadHTMLGlob("view/*.tmpl")
	router.Static("/resource", "./resource")
	router.StaticFile("/favicon.ico", "./resource/favicon.ico")

	router.GET("/", func(ctx *gin.Context) {
		sche := cal.GetScheduleJson()
		cfg := getConfig()

		ctx.HTML(http.StatusOK, "main.tmpl", gin.H{
			"title":   "Top",
			"update":  sche.Update,
			"sche":    getTemplateSche(sche.Schedules),
			"circles": getCircles(),
			"message": cfg.Msg,
		})
	})

	router.GET("/about", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "about.tmpl", gin.H{
			"title": "About",
		})
	})

	router.Run(":52417")

	log.Println("End Introquiz Portal Square Azure")
}

func noescape(tmpl string) template.HTML {
	return template.HTML(tmpl)
}

func loopGet() {
	for {
		callSchedule()
		time.Sleep(time.Minute * 5)
	}
}

func callSchedule() {
	err := cal.MakeScheduleJson()
	if err == nil {
		log.Println("Complete to get events from Google Calendar.")
	} else {
		log.Println("Fail to get events from Google Calendar.")
	}
}

type TmpSchedule struct {
	Schedule   cal.IntroSchedule
	Simple     string
	CircleName string
	Closed     bool
}

func getCircles() map[string]cal.Circle {
	js, err := ioutil.ReadFile("circles.json")
	if err != nil {
		log.Fatalf("Can't read circles.json: %v\n", err)
	}

	var circles map[string]cal.Circle
	err = json.Unmarshal(js, &circles)
	if err != nil {
		log.Fatalf("Unmarshal error circles.json: %v\n", err)
	}

	return circles
}

func getTemplateSche(sche map[string]cal.IntroSchedule) (sc []TmpSchedule) {
	cir := getCircles()
	for _, s := range sche {
		sc = append(sc, TmpSchedule{
			Schedule:   s,
			Simple:     cir[s.CircleId].SimpleName,
			CircleName: cir[s.CircleId].Name,
			Closed:     len(cir[s.CircleId].Overview) == 0,
		})
	}
	sort.Slice(sc, func(i, j int) bool { return sc[i].Schedule.No < sc[j].Schedule.No })
	return sc
}

type Config struct {
	Msg Message `json:"message"`
}

type Message struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func getConfig() (cfg Config) {
	js, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatalf("Can't read config.json: %v\n", err)
	}

	err = json.Unmarshal(js, &cfg)
	if err != nil {
		log.Fatalf("Unmarshal error config.json: %v\n", err)
	}

	return cfg
}
