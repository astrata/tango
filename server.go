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
	"encoding/json"
	"fmt"
	"github.com/astrata/tango/body"
	"github.com/astrata/tango/clf"
	"github.com/astrata/tango/config"
	"github.com/gosexy/to"
	"log"
	"math"
	"net"
	"net/http"
	"net/http/fcgi"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// Server structure, provides Context for every request.
type Server struct {
	serveMux *http.ServeMux
	routes   map[string][]interface{}

	listener net.Listener

	Context *Context
}

// Allocates a new &Server{}.
func NewServer() *Server {
	s := &Server{}

	s.serveMux = http.NewServeMux()
	s.routes = make(map[string][]interface{})

	return s
}

// Starts a fastcgi/http server.
func (server *Server) Run() error {

	var err error

	server.serveMux.Handle("/", server)

	domain := "unix"
	addr := to.String(config.Get("server/socket"))

	if addr == "" {
		domain = "tcp"
		addr = fmt.Sprintf("%s:%d", to.String(config.Get("server/bind")), to.Int(config.Get("server/port")))
	}

	server.listener, err = net.Listen(domain, addr)

	if err != nil {
		log.Fatalf("Failed to bind on %s: %s", addr, err.Error())
	}

	defer server.listener.Close()

	log.Printf("%s is ready to dance.\n", server.listener.Addr())
	log.Printf("Stop server with ^C.\n")

	fmt.Fprintf(os.Stderr, "\n")

	switch to.String(config.Get("server/type")) {
	case "fastcgi":
		if err == nil {
			fcgi.Serve(server.listener, server.serveMux)
		} else {
			log.Fatalf("Failed to start FastCGI server.")
		}
	default:
		if err == nil {
			http.Serve(server.listener, server.serveMux)
		} else {
			log.Fatalf("Failed to start HTTP server.")
		}
	}

	return nil
}

// Maps a route to an interface{}
func (s *Server) Connect(path string, fn interface{}) {
	path = strings.ToLower(path)
	path = fmt.Sprintf("/%s", strings.Trim(path, "/"))

	s.routes[path] = append(s.routes[path], fn)
}

// Routes a *Context to an interface{}
func (server *Server) Route(context *Context) {
	// TODO: Should clean this whole method.

	var name string
	var status int
	var content []byte

	status = 404
	content = []byte{}

	path := strings.ToLower(context.Request.URL.Path)

	path = strings.Trim(path, "/")

	chunks := strings.Split(path, "/")

	// Default content type
	context.SetHeader("Content-Type", "text/html; charset=utf8")

	// Checking for the first chunk that matches a map.
	for i := len(chunks); i >= 0; i-- {

		name = fmt.Sprintf("/%s", strings.Join(chunks[0:i], "/"))

		_, exists := server.routes[name]

		// This map exists.
		if exists {

			var fn interface{}
			var method reflect.Method
			var offset int

			methodExists := false

			for j := 0; j < len(server.routes[name]) && methodExists == false; j++ {

				fn = server.routes[name][j]

				// Routing to Index() if no method was specified
				methodName := "Index"

				// Args offset
				offset = i

				if len(chunks) > i {
					// Translating method name.
					section := chunks[i : i+1][0]
					if section != "" {
						methodName = strings.Trim(section, " ")
						methodName = strings.Replace(methodName, "_", "-", -1)
						methodName = strings.Title(methodName)
						methodName = strings.Replace(methodName, "-", "", -1)
					}
					offset = i + 1
				}

				stype := reflect.TypeOf(fn)

				method, methodExists = stype.MethodByName(methodName)

				if methodExists == false {
					method, methodExists = stype.MethodByName("CatchAll")
					if methodExists == true {
						offset = i
					}
				}
			}

			// Method methodExists
			if methodExists {

				//log.Printf("/%s -> %v", path, method)

				// Copying context into Model.
				if reflect.ValueOf(fn).Elem().FieldByName("Context").IsValid() == true {
					reflect.ValueOf(fn).Elem().FieldByName("Context").Set(reflect.ValueOf(context))
				}

				if reflect.ValueOf(fn).Elem().FieldByName("Params").IsValid() == true {
					reflect.ValueOf(fn).Elem().FieldByName("Params").Set(reflect.ValueOf(context.Params))
				}

				if reflect.ValueOf(fn).Elem().FieldByName("Files").IsValid() == true {
					reflect.ValueOf(fn).Elem().FieldByName("Files").Set(reflect.ValueOf(context.Files))
				}

				// Number of arguments this func requires.
				var argc = method.Type.NumIn()

				// Allocating arguments space
				args := make([]reflect.Value, 1)

				// Adding self reference.
				args[0] = reflect.ValueOf(fn)

				// Starting from offset
				chunkCount := len(chunks)

				a := 1

				// Appending every passed argument.
				for j := offset; len(args) < argc; j++ {

					// Getting current argument type.
					argn := int(math.Min(float64(a), float64(argc-1)))

					currentType := method.Type.In(argn)
					currentValue := reflect.Zero(currentType)

					// Current string value
					if j < len(chunks) {

						vstring := chunks[j]

						// Arguments are strings, may need conversion.
						switch currentType.Kind() {
						case reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8, reflect.Int:
							{
								vint, _ := strconv.Atoi(vstring)
								currentValue = reflect.ValueOf(vint)
							}
						case reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8, reflect.Uint:
							{
								vint, _ := strconv.Atoi(vstring)
								currentValue = reflect.ValueOf(vint)
							}
						case reflect.Float64:
							{
								vfloat64, _ := strconv.ParseFloat(vstring, 64)
								currentValue = reflect.ValueOf(vfloat64)
							}
						case reflect.Float32:
							{
								vfloat32, _ := strconv.ParseFloat(vstring, 32)
								currentValue = reflect.ValueOf(vfloat32)
							}
						case reflect.Bool:
							{
								var vbool bool
								if vstring == "true" || vstring == "1" {
									vbool = true
								} else {
									vbool = false
								}
								currentValue = reflect.ValueOf(vbool)
							}
						case reflect.Slice:
							{
								if method.Type.IsVariadic() == true {
									if a+1 >= argc {
										// Adding all remaining chunks.
										for ; j < chunkCount; j++ {
											currentValue = reflect.Append(currentValue, cast(currentType, chunks[j]))
										}
									} else {
										currentValue = cast(currentType, chunks[j])
									}
								} else {
									panic("Array values are not yet supported.")
								}
							}
						default:
							{
								currentValue = reflect.ValueOf(vstring)
							}
						}
					}

					args = append(args, currentValue)

					a++
				}

				// Executing called method.
				var output []reflect.Value

				if method.Type.IsVariadic() == true {
					output = method.Func.CallSlice(args)
				} else {
					output = method.Func.Call(args)
				}

				// Callback after execution.
				context.afterExecute()

				status = 200
				context.Writer.Header().Set("Content-type", "text/plain; charset=utf8")

				if len(output) > 0 {

					value := output[0].Interface()

					switch output[0].Interface().(type) {
					case body.Body:

						for k, v := range value.(body.Body).Header() {
							context.Writer.Header()[k] = v
						}

						content = value.(body.Body).Get()

						status = value.(body.Body).Status()
					case string:
						context.Writer.Header().Set("Content-type", "text/html; charset=utf8")
						content = []byte(value.(string))
					case nil:
						status = 404
					default:
						var err error
						content, err = json.Marshal(value)
						if err == nil {
							context.Writer.Header().Set("Content-type", "text/plain; charset=utf8")
						} else {
							status = 500
						}
					}

				} else {
					status = 404
					content = []byte{}
				}

			}

			break
		}
	}

	size := len(content)

	if status != 200 {
		http.Error(context.Writer, http.StatusText(status), status)
	}

	context.Writer.Write(content)

	clf.Print(context.Request, status, size)
}

// Interface method for handling HTTP.
func (server *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	context := newContext(server, writer, request)
	server.Route(context)
}
