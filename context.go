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
	"github.com/astrata/tango/config"
	"github.com/gosexy/to"
	"net/http"
	"reflect"
	"strconv"
)

// Each request has its own Context struct that contains info about the request type and provides
// methods for reading parameters.
type Context struct {

	// Request type
	GET    bool
	POST   bool
	PUT    bool
	DELETE bool

	// Standard response writer
	Writer http.ResponseWriter

	// Standard request
	Request *http.Request

	// Server
	Server *Server

	Params Value

	Cookies Value

	Files Files

	cookieMap map[string]*http.Cookie
}

func newContext(server *Server, writer http.ResponseWriter, request *http.Request) *Context {

	context := &Context{}

	switch request.Method {
	case "GET":
		context.GET = true
	case "POST":
		context.POST = true
	case "DELETE":
		context.DELETE = true
	case "PUT":
		context.PUT = true
	}

	context.cookieMap = make(map[string]*http.Cookie)

	context.Server = server
	context.Request = request
	context.Writer = writer

	maxSize := to.Int64(config.Get("server/request_max_size"))

	request.ParseMultipartForm(maxSize)

	context.Params = context.getParams()
	context.Cookies = context.getCookies()
	context.Files = context.getFiles()

	server.Context = context

	return context
}

func cast(t reflect.Type, value string) reflect.Value {
	result := reflect.Zero(t)
	// Is there a cleaner way of doing this?
	switch t.String() {
	case "[]int64", "[]int32", "[]int16", "[]int8", "[]int":
		{
			vint, _ := strconv.Atoi(value)
			result = reflect.ValueOf(vint)
		}
	case "[]string":
		{
			result = reflect.ValueOf(value)
		}
	}
	return result
}

// Creates a cookie for the current session.
func (context *Context) Cookie(name string) *http.Cookie {
	cookie := &http.Cookie{}
	cookie.Name = name

	context.cookieMap[name] = cookie

	return cookie
}

// Sends a HTTP error code.
func (context *Context) HttpError(code int) {
	http.Error(context.Writer, http.StatusText(code), code)
}

// Sends a 301 Location header with the given value.
func (context *Context) Redirect(value string) {
	context.Writer.Header().Set("Location", value)
	context.HttpError(301)
}

// Sets a header to be sent to the client.
func (context *Context) SetHeader(name string, value string) {
	context.Writer.Header().Set(name, value)
}

func (context *Context) afterExecute() {
	for key, _ := range context.cookieMap {
		http.SetCookie(context.Writer, context.cookieMap[key])
	}
}

func (context *Context) getFiles() Files {
	params := Files{}

	if context.Request.MultipartForm != nil {
		for key, val := range context.Request.MultipartForm.File {
			params.Set(key, val)
		}
	}

	return params
}

func (context *Context) getCookies() map[string]interface{} {
	params := Value{}

	cookies := context.Request.Cookies()

	itop := len(cookies)

	for i := 0; i < itop; i++ {
		params[cookies[i].Name] = cookies[i].Value
	}

	return params
}

func (context *Context) getParams() Value {

	params := Value{}

	for key, val := range context.Request.Form {
		params.Set(key, val)
	}

	if context.Request.MultipartForm != nil {
		for key, val := range context.Request.MultipartForm.Value {
			params.Set(key, val)
		}
	}

	return params
}
