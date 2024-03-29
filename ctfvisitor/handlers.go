package ctfvisitor

import (
	"context"
	"log"
	"time"

	"github.com/tebeka/selenium"
)

type Handler func(ctx context.Context, wd selenium.WebDriver) error

// VisitHandler is a generic wrapper around Get() that also installs cookies
// first before visiting the page
func CookieHandler(
	ctx context.Context, wd selenium.WebDriver,
	path string, cookies []*selenium.Cookie) error {
	if err := wd.SetImplicitWaitTimeout(200); err != nil {
		return err
	}

	wd.Get(path)
	for _, cookie := range cookies {
		log.Printf("visiting %v with cookie %v\n", path, cookie)
		if err := wd.AddCookie(cookie); err != nil {
			return err
		}
	}

	if err := wd.Get(path); err != nil {
		return err
	}

	time.Sleep(1000 * time.Millisecond)

	return nil

}
