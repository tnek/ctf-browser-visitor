package ctfvisitor

import (
	"context"
	"log"

	"github.com/tebeka/selenium"
)

type Handler func(ctx context.Context, wd selenium.WebDriver) error

// VisitHandler is a generic wrapper around Get() that also installs cookies
// first before visiting the page
func CookieHandler(
	ctx context.Context, wd selenium.WebDriver,
	path string, cookies []*selenium.Cookie) error {
	log.Printf("visiting %v with cookies %v\n", path, cookies)
	for _, cookie := range cookies {
		if err := wd.AddCookie(cookie); err != nil {
			return err
		}
	}

	return wd.Get(path)

}
