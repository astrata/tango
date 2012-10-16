package models

import (
	"fmt"
	"github.com/astrata/tango"
	"github.com/astrata/tango/app"
	"github.com/astrata/tango/body"
	"os"
	"strings"
)

// Root directory for static files.
var Root = "static"

type Static struct {
}

func init() {
	app.Register("Static", &Static{})
	app.Fallback("/", app.App("Static"))
}

// Checking root directory on start.
func (self *Static) StartUp() {
	info, err := os.Stat(Root)
	if err == nil {
		if info.IsDir() == false {
			panic(fmt.Sprintf("%s is not a directory.\n", Root))
		}
	} else {
		panic(err.Error())
	}
}

// Catches all requests and serves files.
func (self *Static) CatchAll(path ...string) body.Body {

	content := body.File()

	filename := Root + tango.PS + strings.Trim(strings.Join(path, tango.PS), tango.PS)

	info, err := os.Stat(filename)

	if err == nil {

		if info.IsDir() == true {

			filename = strings.Trim(filename, tango.PS) + tango.PS + "index.html"

			info, err = os.Stat(filename)

			if err != nil {
				return nil
			}

			if info.IsDir() == true {
				return nil
			}

		}

		content.Set(filename)

		return content
	}

	return nil
}
