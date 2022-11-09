package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/tebeka/selenium"
	"github.com/tnek/ctf-browser-visitor/ctfvisitor"
)

type Job struct {
	URL     string             `json:"url"`
	Cookies []*selenium.Cookie `json:"cookies"`
}

type App struct {
	Ctf *ctfvisitor.Dispatch
}

func (a *App) Index(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("job")

	log.Printf("deserializing %v", q)
	var job Job
	if err := json.Unmarshal([]byte(q), &job); err != nil {
		log.Printf("failed to unmarshal: %v", err)
	}
	log.Printf("job json: %v\n", job)
	site := &ctfvisitor.Site{Path: job.URL, Cookies: job.Cookies}

	if err := a.Ctf.Queue(site); err != nil {
		log.Printf("error with handling request '%v': %v", q, err)
	}
}

func man() {
	fmt.Printf("%v: [--selenium <selenium-server.jar path>] [--driver <[chrome|gecko]driver path>] [-q <queue size>] [-w <number of workers>] host port\n", os.Args[0])
}

func main() {
	SELENIUM_PATH := flag.String("selenium", "./selenium-server.jar", "Path to the selenium-server.jar.")
	BROWSER_PATH := flag.String("driver", "/usr/bin/chromedriver", "Path to the [gecko|chrome]driver binary.")
	QUEUE_SIZE := flag.Int("q", 10000, "Max number of queued up requests")
	NUMWORKERS := flag.Int("w", 10, "Number of instances of drivers")

	flag.Parse()

	if flag.NArg() < 2 {
		man()
		return
	}

	host := flag.Arg(0)
	port, err := strconv.Atoi(flag.Arg(1))
	if err != nil {
		fmt.Printf("Invalid argument for port: %v", flag.Arg(1))
		os.Exit(1)
	}

	browserType := ctfvisitor.UNKNOWN
	if strings.Contains((*BROWSER_PATH), "chromedriver") {
		browserType = ctfvisitor.CHROME
	} else if strings.Contains((*BROWSER_PATH), "geckodriver") {
		browserType = ctfvisitor.FIREFOX
	}

	cfg := &ctfvisitor.Config{
		QueueSize:    *QUEUE_SIZE,
		SeleniumPath: *SELENIUM_PATH,
		BrowserPath:  *BROWSER_PATH,
		Browser:      browserType,
		MinPort:      port + 1,
	}

	ctf, err := ctfvisitor.Init(cfg)
	if err != nil {
		log.Panicf("failed to initialize: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go ctf.LoopWithRestart(ctx, *NUMWORKERS)

	a := &App{
		Ctf: ctf,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/visit", a.Index)

	s := http.Server{
		Addr:    fmt.Sprintf("%v:%v", host, port),
		Handler: mux,
	}

	log.Fatal(s.ListenAndServe())
}
