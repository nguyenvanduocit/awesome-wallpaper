package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/mileusna/crontab"
	"github.com/olekukonko/tablewriter"
	"github.com/reujab/wallpaper"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"time"
)

const baseURL = "https://source.unsplash.com/random"
const version = "1.0.0"

type Schedule struct {
	Description string `json:"description,omitempty"`
	Schedule    string `json:"schedule,omitempty"`
	Keywords    string `json:"keywords,omitempty"`
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	var table *tablewriter.Table
	table = tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Schedule", "Keyword", "Description"})
	var (
		schedule       string
		keywords       string
		configFilePath string
		showVersion    bool
		showHelp       bool
	)

	flag.BoolVar(&showVersion, "version", false, fmt.Sprintf("Current version: %s", version))
	flag.BoolVar(&showHelp, "help", false, "View help message")
	flag.StringVar(&schedule, "schedule", "30 * * * *", "(optional) A crontab-like syntax schedule")
	flag.StringVar(&keywords, "keywords", "", "(optional) Keyword to search for image")
	flag.StringVar(&configFilePath, "conf", "", "(optional) Config file path")
	flag.Parse()

	if showHelp {
		flag.PrintDefaults()
		os.Exit(0)
	}

	if showVersion {
		log.Fatalf("awesome-wallpaper %s\n", version)
		os.Exit(0)
	}

	var schedules []Schedule

	if configFilePath != "" {
		var err error
		schedules, err = parseScheduleConfig(configFilePath)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		schedules = []Schedule{
			{
				Schedule: schedule,
				Keywords: keywords,
			},
		}
	}

	ctab := crontab.New()
	for _, job := range schedules {
		ctab.MustAddJob(job.Schedule, changeWallpaper, job)
		table.Append([]string{job.Schedule, job.Keywords, job.Description})
	}
	table.Render()
	log.Println("Running...")
	select {}
}

func changeWallpaper(schedule Schedule) error {
	sig := time.Now().Unix()
	url := fmt.Sprintf("%s?sig=%d", baseURL, sig)
	if schedule.Keywords != "" {
		url = fmt.Sprintf("%s&%s", url, schedule.Keywords)
	}
	imagePath, err := downloadImage(url)
	if err != nil {
		return err
	}
	if err := wallpaper.SetFromFile(imagePath); err != nil {
		return err
	}
	return nil
}

func parseScheduleConfig(configFilePath string) ([]Schedule, error) {
	plan, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}
	var schedules []Schedule
	if err := json.Unmarshal(plan, &schedules); err != nil {
		return nil, err
	}

	return schedules, nil
}

func downloadImage(url string) (string, error) {
	tmpFile, err := openTempFile()
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return "", errors.New("non-200 status code")
	}

	_, err = io.Copy(tmpFile, res.Body)
	if err != nil {
		return "", err
	}

	return tmpFile.Name(), nil
}

func openTempFile() (*os.File, error) {
	filename, err := ensureCacheFile()
	if err != nil {
		return nil, err
	}
	tmpFile, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0600)
	return tmpFile, err
}

func ensureCacheFile() (string, error) {
	_, caller, _, _ := runtime.Caller(1)
	dir := path.Dir(caller)
	files, err := filepath.Glob(fmt.Sprintf("%s/background-*.jpg", dir))
	if err != nil {
		return "", err
	}
	for _, file := range files {
		if err := os.Remove(file); err != nil {
			return "", err
		}
	}
	filename := fmt.Sprintf("%s/background-%d.jpg", path.Dir(caller), time.Now().Unix())
	return filename, nil
}
