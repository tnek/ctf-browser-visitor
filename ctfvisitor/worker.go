package ctfvisitor

import (
	"context"
	"fmt"
	"log"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"github.com/tebeka/selenium/firefox"
)

type Worker struct {
	ID int

	wd  selenium.WebDriver
	srv *selenium.Service
}

type WorkerConfig struct {
	Port int
	// Dest must be a format string that includes Port.
	// ex: "http://localhost:%d/wd/hub"
	Dest string

	// SeleniumPath is the Selenium .jar.
	SeleniumPath string

	Caps        selenium.Capabilities
	ServiceOpts []selenium.ServiceOption
}

// DefaultWC gives reasonable defaults as selenium configs per browser for a
// worker.
func DefaultWC(browser BrowserType, path string) *WorkerConfig {
	switch browser {
	case UNKNOWN:
		fallthrough
	case CHROME:
		caps := selenium.Capabilities{"browserName": "chrome"}
		caps.AddChrome(chrome.Capabilities{
			Args: []string{"--no-sandbox", "--headless"},
		})

		return &WorkerConfig{
			Dest: "http://localhost:%v/wd/hub",
			Caps: caps,
			ServiceOpts: []selenium.ServiceOption{
				selenium.ChromeDriver(path),
			},
		}
	case FIREFOX:
		caps := selenium.Capabilities{"browserName": "firefox"}
		caps.AddFirefox(firefox.Capabilities{
			Args: []string{"-headless"},
		})

		return &WorkerConfig{
			Dest: "http://localhost:%v/wd/hub",
			Caps: caps,
			ServiceOpts: []selenium.ServiceOption{
				selenium.GeckoDriver(path),
			},
		}
	}
	return nil
}

func InitWorker(wc *WorkerConfig, id int, port int) (*Worker, func(), error) {
	srv, err := selenium.NewSeleniumService(wc.SeleniumPath, port, wc.ServiceOpts...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to initialize selenium: %w", err)
	}

	wd, err := selenium.NewRemote(wc.Caps, fmt.Sprintf(wc.Dest, port))
	if err != nil {
		return nil, func() { srv.Stop() }, fmt.Errorf("failed to initialize selenium: %w", err)
	}

	return &Worker{ID: id, srv: srv, wd: wd},
		func() {
			wd.Quit()
			srv.Stop()
		}, nil
}

func (w *Worker) Run(ctx context.Context, wq chan Handler) error {
	for {
		select {
		case handle := <-wq:
			if err := handle(ctx, w.wd); err != nil {
				log.Printf("handler failed with error: %v", err)
			}
			if err := w.Reset(ctx); err != nil {
				return fmt.Errorf("cleanup of worker failed: %w", err)
			}

		case <-ctx.Done():
			if err := w.Cleanup(ctx); err != nil {
				return fmt.Errorf("cleanup of worker failed: %w", err)
			}
			return nil
		}
	}
	return nil
}

func (w *Worker) Reset(ctx context.Context) error {
	if err := w.wd.DeleteAllCookies(); err != nil {
		return fmt.Errorf("reset failed to clear cookies: %w", err)
	}
	return nil
}

func (w *Worker) Cleanup(ctx context.Context) error {
	return nil
}
