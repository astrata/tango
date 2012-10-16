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
	"fmt"
	"mime"
	"net/http"
	"os"
	"path/filepath"
)

type fileContent struct {
	header  http.Header
	status  int
	content []byte

	ForceDownload bool

	Name     string
	MIMEType string
}

// Returns a Body that can be used to send files to the client.
func File() Body {
	self := &fileContent{}
	self.status = 200
	self.header = http.Header{}
	self.MIMEType = "application/octect-stream"
	self.Name = "file.bin"
	self.ForceDownload = false
	return self
}

// Returns the headers to be sent along the response.
func (self *fileContent) Header() http.Header {
	self.header.Add("Content-type", self.MIMEType)

	if self.ForceDownload == true {
		self.header.Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s;", filepath.Base(self.Name)))
	}

	return self.header
}

// Uses the given argument to create a File.
// Currently it only accepts a filename, as a string.
func (self *fileContent) Set(value interface{}) {
	switch value.(type) {
	case string:
		self.loadFromDisk(value.(string))
	default:
		self.status = 404
	}
}

func (self *fileContent) loadFromDisk(filename string) {

	info, err := os.Stat(filename)

	if err == nil {

		if info.IsDir() == true {
			self.status = 404
			return
		}

		file, err := os.Open(filename)

		if err != nil {
			self.status = 500
			return
		}

		defer file.Close()

		self.content = make([]byte, info.Size())

		file.Read(self.content)

		self.Name = filepath.Base(filename)

		self.MIMEType = mime.TypeByExtension(filepath.Ext(filename))

	}
}

// Returns the file contents.
func (self *fileContent) Get() []byte {
	return self.content
}

// Returns the request HTTP status.
func (self *fileContent) Status() int {
	return self.status
}
