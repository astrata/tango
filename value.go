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

package tango

import (
	"fmt"
	"mime/multipart"
	"net/textproto"
	//"strconv"
	//"strings"
)

// This struct holds multipart files.
type File struct {
	header *multipart.FileHeader
	value  string

	Name   string
	Header textproto.MIMEHeader
}

// Returns the contents of a File.
func (file File) Content() []byte {

	var buf []byte

	fp, err := file.header.Open()

	if err == nil {
		defer fp.Close()

		end, _ := fp.Seek(0, 2)
		fp.Seek(0, 0)

		buf = make([]byte, end)

		fp.Read(buf)
	}

	return buf
}

// A file map
type Files map[string]interface{}

// Adds an array of files to the map.
func (files Files) Set(name string, v []*multipart.FileHeader) {
	files[name] = v
}

// Returns the first file that matches the given name.
func (files Files) Get(name string) *File {
	header := files.GetAll(name)

	if header == nil {
		return nil
	}

	return header[0]
}

// Returns all files that match the given name.
func (files Files) GetAll(name string) []*File {
	res := []*File{}

	if files[name] == nil {
		return nil
	}

	all := files[name].([]*multipart.FileHeader)

	if all == nil {
		return nil
	}

	for i := 0; i < len(all); i++ {
		file := &File{header: all[i]}
		file.Name = file.header.Filename
		file.Header = file.header.Header
		res = append(res, file)
	}

	return res
}

// For storing and retrieving POST and GET values.
type Value map[string]interface{}

// Sets a value.
func (value Value) Set(name string, v interface{}) {
	value[name] = v
}

// Returns only the given keys in a new Value.
func (value Value) Filter(names ...string) Value {
	response := Value{}

	for _, key := range names {
		if _, exists := value[key]; exists {
			response[key] = value[key]
		}
	}

	return response
}

// Returns all the associated values.
func (value Value) GetAll(name string) interface{} {
	return value[name]
}

// Returns the first value associated with a name, as a string.
func (value Value) Get(name string) string {
	v := value.GetAll(name)
	if v == nil {
		return ""
	}
	return value[name].([]string)[0]
}

/*
func (value Value) GetList(name string) []string {
	v := value[name]
	if v == nil {
		return []string{}
	}
	return v.([]string)
}

func (value Value) GetString(name string) string {
	v := value[name]
	if v == nil {
		return ""
	}
	return v.([]string)[len(v.([]string))-1]
}
*/

// Verifies that all the given names have a value.
func (value Value) Require(name ...string) (bool, map[string][]string) {
	messages := make(map[string][]string)

	valid := true

	for _, key := range name {
		if value.Get(key) == "" {
			valid = false
			if messages[key] == nil {
				messages[key] = []string{}
			}
			messages[key] = append(messages[key], fmt.Sprintf("Missing required parameter %s.", key))
		}
	}

	return valid, messages
}

/*
func (value Value) GetInt(name string) int {
	i, _ := strconv.Atoi(value.GetString(name))
	return i
}

func (value Value) GetFloat(name string) float64 {
	f, _ := strconv.ParseFloat(value.GetString(name), 64)
	return f
}

func (value Value) GetBool(name string) bool {
	b := strings.ToLower(value.GetString(name))
	if b == "0" || b == "false" || b == "" {
		return false
	}
	return true
}
*/
