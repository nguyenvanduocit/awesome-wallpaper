package main

import (
	"flag"
	"fmt"
	"github.com/mileusna/crontab"
	"github.com/reujab/wallpaper"
	"log"
	"time"
)

func changeWallpaper(keywords string) {
	sig := time.Now().Unix()
	url := fmt.Sprintf("https://source.unsplash.com/random?sig=%d", sig)
	if keywords != "" {
		url = fmt.Sprintf("%s&%s", url, keywords)
	}
	if err := wallpaper.SetFromURL(url); err != nil {
		log.Println(err)
	}
}

func main() {
	var (
		schedule string
		keywords string
	)
	flag.StringVar(&schedule, "schedule", "* * * * *", "Crontab-ilke schedule string")
	flag.StringVar(&keywords, "keywords", "", "Keyword to search for image")
	flag.Parse()
	ctab := crontab.New()
	ctab.MustAddJob(schedule, changeWallpaper, keywords)
	ctab.RunAll()
	select {}
}
