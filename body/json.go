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

package body

import (
	"encoding/json"
	"github.com/gosexy/sugar"
	"log"
	"net/http"
)

type jsonContent struct {
	status  int
	header  http.Header
	content []byte
	data    interface{}
}

// Returns a Body for an JSON response.
func Json() Body {
	self := &jsonContent{}
	self.status = 200
	self.header = http.Header{}
	self.header.Add("Content-type", "text/plain; charset=utf8")
	return self
}

// Returns the headers to be sent along the request.
func (self *jsonContent) Header() http.Header {

	switch self.data.(type) {
	case map[string]interface{}:
		if self.data.(map[string]interface{})["error"] != nil {
			self.status = 400
		}
	case sugar.Map:
		if self.data.(sugar.Map)["error"] != nil {
			self.status = 400
		}
	}

	return self.header
}

// Returns the request HTTP status.
func (self *jsonContent) Status() int {
	return self.status
}

// Sets the request contents.
func (self *jsonContent) Set(value interface{}) {
	data, err := json.Marshal(value)
	if err == nil {
		self.content = data
		self.data = value
	} else {
		log.Printf("body.Json: %s\n", err)
	}
}

// Returns the request contents that are going to be written.
func (self *jsonContent) Get() []byte {
	return self.content
}
