package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/mileusna/crontab"
	"github.com/olekukonko/tablewriter"
	"github.com/reujab/wallpaper"
	"github.com/takama/daemon"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"
)

const baseURL = "https://source.unsplash.com/random"
const version = "1.0.0"

type Service struct {
	daemon.Daemon
}
type Schedule struct {
	Description string `json:"description,omitempty"`
	Schedule    string `json:"schedule,omitempty"`
	Keywords    string `json:"keywords,omitempty"`
}

var (
	name        = "vn.12bit.awesome-wallpaper"
	description = "awesome-wallpaper"
)

var (
	schedule       string
	keywords       string
	configFilePath string
	serviceAction  string
	showVersion    bool
	showHelp       bool
	isDeamon       bool
)

func main() {
	flag.BoolVar(&showVersion, "version", false, fmt.Sprintf("Current version: %s", version))
	flag.BoolVar(&showHelp, "help", false, "View help message")
	flag.StringVar(&schedule, "schedule", "30 * * * *", "(optional) A crontab-like syntax schedule")
	flag.StringVar(&keywords, "keywords", "", "(optional) Keyword to search for image")
	flag.StringVar(&configFilePath, "conf", "", "(optional) Config file path")
	flag.StringVar(&serviceAction, "service", "", "(optional) Action about services: install, uninstall, remove, stop, status")
	flag.BoolVar(&isDeamon, "deamon", false, "(optional) Indicate if program is running as deamon")
	flag.Parse()

	if showHelp {
		flag.PrintDefaults()
		os.Exit(0)
	}

	if showVersion {
		log.Printf("awesome-wallpaper %s\n", version)
		os.Exit(0)
	}

	if serviceAction != "" {
		srv, err := daemon.New(name, description)
		if err != nil {
			log.Println("Error: ", err)
			os.Exit(1)
		}
		service := &Service{srv}
		switch serviceAction {
		case "install":
			args := []string{
				"--deamon",
			}
			for _, arg := range os.Args[1:] {
				if strings.Index(arg, "--service=") != 0 {
					args = append(args, arg)
				}
			}
			status, err := service.Install(args...)
			if err != nil {
				log.Fatalln(err)
			}
			log.Println(status)
			log.Println(args)
			os.Exit(0)
		case "remove":
			status, err := service.Remove()
			if err != nil {
				log.Fatalln(err)
			}
			log.Println(status)
			os.Exit(0)
		case "start":
			status, err := service.Start()
			if err != nil {
				log.Fatalln(err)
			}
			log.Println(status)
			os.Exit(0)
		case "stop":
			status, err := service.Stop()
			if err != nil {
				log.Fatalln(err)
			}
			log.Println(status)
			os.Exit(0)
		case "status":
			status, err := service.Status()
			if err != nil {
				log.Fatalln(err)
			}
			log.Println(status)
			os.Exit(0)
		default:
			log.Println("Usage: awesome-wallpaper service install | remove | start | stop | status")
			os.Exit(0)
		}
	}

	setupLogger()

	var schedules []Schedule
	if configFilePath != "" {
		var err error
		if schedules, err = parseScheduleConfig(configFilePath); err != nil {
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

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Schedule", "Keyword", "Description"})

	ctab := crontab.New()
	for _, job := range schedules {
		ctab.MustAddJob(job.Schedule, changeWallpaper, job)
		table.Append([]string{job.Schedule, job.Keywords, job.Description})
	}
	table.Render()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)

	log.Println("Running...")

	killSignal := <-interrupt
	log.Println("Got signal:", killSignal)
	if killSignal == os.Interrupt {
		log.Println("Interruped by system signal ")
	}
	log.Println("Bye...")
}

func changeWallpaper(schedule Schedule) error {
	sig := time.Now().Unix()
	url := fmt.Sprintf("%s?sig=%d", baseURL, sig)
	if schedule.Keywords != "" {
		url = fmt.Sprintf("%s&%s", url, schedule.Keywords)
	}
	imagePath, err := downloadImage(url)
	if err != nil {
		log.Println(err)
		return err
	}
	if err := wallpaper.SetFromFile(imagePath); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func parseScheduleConfig(configFilePath string) ([]Schedule, error) {
	if !filepath.IsAbs(configFilePath) {
		_, caller, _, _ := runtime.Caller(1)
		dir := path.Dir(caller)
		configFilePath = filepath.Join(dir, configFilePath)
	}
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
	tmpFile, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	return tmpFile, err
}

func ensureCacheFile() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", err
	}
	dir := filepath.Dir(ex)
	files, err := filepath.Glob(fmt.Sprintf("%s/background-*.jpg", dir))
	if err != nil {
		return "", err
	}
	for _, file := range files {
		if err := os.Remove(file); err != nil {
			return "", err
		}
	}
	filename := fmt.Sprintf("%s/background-%d.jpg", dir, time.Now().Unix())
	return filename, nil
}
