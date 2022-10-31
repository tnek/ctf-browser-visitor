package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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

//http://127.0.0.1:5000/visit?job={"url":"url_to_visit","cookies":{"key":"value"}}
func (a *App) Index(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("job")

	var job Job
	json.Unmarshal([]byte(q), &job)

	handler := ctfvisitor.CookieHandler(job.URL, job.Cookies)
	if err := a.Ctf.Queue(handler); err != nil {
		log.Printf("error with handling request '%v': %v", q, err)
	}
}

func main() {
	const (
		numWorkers = 10
		host       = "0.0.0.0"
		port       = 8080
	)

	cfg := &ctfvisitor.Config{
		QueueSize:    1000,
		SeleniumPath: "./selenium-server.jar",
		BrowserPath:  "./chromedriver",
		Browser:      ctfvisitor.Chrome,
	}

	ctf, err := ctfvisitor.Init(cfg)
	if err != nil {
		log.Panicf("%v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go ctf.LoopWithRestart(ctx, numWorkers)

	a := &App{
		Ctf: ctf,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", a.Index)

	s := http.Server{
		Addr:    fmt.Sprintf("%v:%v", host, port),
		Handler: mux,
	}

	log.Fatal(s.ListenAndServe())

}
