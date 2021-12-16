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

	"github.com/pgDora56/Azure/cal"
	"github.com/pgDora56/Azure/dynamodb"
	"github.com/pgDora56/Azure/schemas"
)

func main() {
	if len(os.Args) > 1 {
		if os.Args[1] == "--get" {
			cfg := getConfig()
			callSchedule(cfg.Dynamo)
			return
		}
		if os.Args[1] == "--show" {
			cfg := getConfig()
			data, err := dynamodb.Get(cfg.Dynamo)
			log.Println(data)
			if err != nil {
				panic(err)
			}
			return
		}
		log.Println("Unknown arguments")
		return
	}

	log.Println("Start Introquiz Portal Square Azure")

	cfg := getConfig()
	go loopGet(cfg.Dynamo, cfg.CheckDistance)

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

func loopGet(dcfg schemas.DynamoConfig, distance int) {
	for {
		callSchedule(dcfg)
		time.Sleep(time.Minute * time.Duration(distance))
	}
}

func callSchedule(dcfg schemas.DynamoConfig) {
	err := cal.Insert2Dynamo(dcfg)
	if err == nil {
		log.Println("Complete to get events from Google Calendar.")
	} else {
		log.Println("Fail to get events from Google Calendar.")
		log.Println(err)
	}
}

type TmpSchedule struct {
	Schedule   schemas.IntroSchedule
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

func getTemplateSche(sche map[string]schemas.IntroSchedule) (sc []TmpSchedule) {
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

func getConfig() (cfg schemas.Config) {
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
