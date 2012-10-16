/*
  Tango!

  Copyright (c) 2012 Astrata Software, <http://astrata.mx>
  Written by José Carlos Nieto <xiam@menteslibres.org>

  Permission is hereby granted, free of charge, to any person obtaining
  a copy of this software and associated documentation files (the
  "Software"), to deal in the Software without restriction, including
  without limitation the rights to use, copy, modify, merge, publish,
  distribute, sublicense, and/or sell copies of the Software, and to
  permit persons to whom the Software is furnished to do so, subject to
  the following conditions:

  The above copyright notice and this permission notice shall be
  included in all copies or substantial portions of the Software.

  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
  EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
  MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
  NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
  LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
  OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
  WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

/*
  This package implements the Common Log Format for any *http.Requests and
  prints the result to os.Stdout.

  [1]: http://en.wikipedia.org/wiki/Common_Log_Format
*/
package clf

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

func chunk(value string) string {
	if value == "" {
		return "-"
	}
	return value
}

// Prints a Common Log Format of a request to stdout.
func Print(req *http.Request, status int, size int) {
	fmt.Println(
		strings.Join([]string{
			chunk(req.RemoteAddr),
			chunk(""), // TODO: http://tools.ietf.org/html/rfc1413 IDENT
			chunk(""), // TODO: User id
			chunk("[" + time.Now().Format("02/Jan/2006:15:04:05 -0700") + "]"),
			chunk("\"" + fmt.Sprintf("%s %s %s", req.Method, req.RequestURI, req.Proto) + "\""),
			chunk(fmt.Sprintf("%d", status)),
			chunk(fmt.Sprintf("%d", size)),
		},
			" ",
		),
	)
}
