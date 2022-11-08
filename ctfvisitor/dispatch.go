package ctfvisitor

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"

	"github.com/tebeka/selenium"
)

type BrowserType int

const (
	UNKNOWN BrowserType = iota
	CHROME
	FIREFOX
)

type Site struct {
	Path    string
	Cookies []*selenium.Cookie
}

type Config struct {
	// SeleniumPath is the Selenium .jar.
	SeleniumPath string

	// QueueSize is the maximum number of handler events we can queue up at
	// once.
	QueueSize int

	Browser     BrowserType
	BrowserPath string

	MinPort int
	MaxPort int
}

type Dispatch struct {
	MinPort int
	MaxPort int

	wc        *WorkerConfig
	wq        chan *Site
	assignLck sync.Mutex
	idToPort  map[int]int
	portToId  map[int]int
}

func Init(c *Config) (*Dispatch, error) {
	return InitWithWC(c, DefaultWC(c.Browser, c.BrowserPath))
}

func InitWithWC(c *Config, wc *WorkerConfig) (*Dispatch, error) {
	if c.MinPort == 0 {
		c.MinPort = 1001
	}
	if c.MinPort <= 1000 {
		return nil, fmt.Errorf("attempting to use reserved port '%d' <= 1000", c.MinPort)
	}
	if c.MaxPort > 65535 || c.MaxPort == 0 {
		c.MaxPort = 65535
	}

	if c.SeleniumPath == "" {
		return nil, fmt.Errorf("missing SeleniumPath in config")
	}

	if c.QueueSize == 0 {
		c.QueueSize = 100000
	}

	wc.SeleniumPath = c.SeleniumPath

	d := &Dispatch{
		MinPort: c.MinPort,
		MaxPort: c.MaxPort,

		wc:       wc,
		wq:       make(chan *Site, c.QueueSize),
		idToPort: map[int]int{},
		portToId: map[int]int{},
	}

	return d, nil
}

func (d *Dispatch) Queue(s *Site) error {
	d.wq <- s
	return nil
}

// assign creates a valid non-conflicting id<->port pairing for a new worker.
func (d *Dispatch) assign() (id int, port int, cleanup func()) {
	d.assignLck.Lock()
	id = rand.Int()
	for _, ok := d.idToPort[id]; ok; {
		id = rand.Int()
	}

	port = rand.Intn(d.MaxPort-d.MinPort) + d.MinPort
	for _, ok := d.portToId[port]; ok; {
		port = rand.Intn(d.MaxPort-d.MinPort) + d.MinPort
	}

	d.idToPort[id] = port
	d.portToId[port] = id

	cleanup = func() {
		d.assignLck.Lock()
		delete(d.portToId, port)
		delete(d.idToPort, id)
		d.assignLck.Unlock()
	}
	d.assignLck.Unlock()
	return
}

// LoopWithRestart maintains a constant worker pool of `workerCount` workers.
func (d *Dispatch) LoopWithRestart(ctx context.Context, workerCount int) {
	tokens := make(chan bool, workerCount)
	for i := 0; i < workerCount; i++ {
		tokens <- true
	}

	for {
		select {
		case <-tokens:
			log.Printf("taking job")
			go func() {
				defer func() { tokens <- true }()

				id, port, portCleanup := d.assign()
				defer portCleanup()

				w, workerCleanup, err := InitWorker(d.wc, id, port)
				if err != nil {
					log.Printf("worker %d failed to initialize: %v", id, err)
					return
				}
				defer workerCleanup()

				if err := w.Run(ctx, d.wq); err != nil {
					log.Printf("worker %d failed: %v", id, err)
					return
				}
			}()
		}
	}
}
