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

package app

import (
	"fmt"
	"github.com/astrata/tango"
	"log"
	"os"
)

var routes = make(map[string]Model)
var fallbacks = make(map[string]Model)
var apps = make(map[string]Model)

// Tango! server.
var Server *tango.Server

// Basic model structure.
// Any model must have, at least, the StartUp() function.
type Model interface {
	StartUp()
}

func init() {
	log.Println("Tango! by Astrata")
	fmt.Fprintf(os.Stderr, "\n")
}

// Registers a Model struct so it can be retrieved by the given name.
func Register(name string, app Model) {
	if _, ok := apps[name]; ok {
		panic(fmt.Sprintf("App %s was already registered.", name))
	}
	apps[name] = app
}

// Returns a previously defined Model struct.
func App(name string) Model {
	if _, ok := apps[name]; ok {
		return apps[name]
	}
	return nil
}

// Defines the main route for a Model.
func Route(name string, app Model) {
	if _, ok := routes[name]; ok {
		panic(fmt.Sprintf("Route %s was already registered.", name))
	}
	routes[name] = app
}

// Like Route() but called only as the last option.
func Fallback(name string, app Model) {
	if _, ok := fallbacks[name]; ok {
		panic(fmt.Sprintf("Fallback %s was already registered.", name))
	}
	fallbacks[name] = app
}

// Initializes a fastcgi/http server.
func Run() {

	log.Println("Initializing server...")

	Server = tango.NewServer()

	for route, model := range routes {
		log.Printf("Adding route: %s\n", route)
		model.StartUp()
		Server.Connect(route, model)
	}

	for fallback, model := range fallbacks {
		log.Printf("Adding fallback: %s\n", fallback)
		model.StartUp()
		Server.Connect(fallback, model)
	}

	fmt.Fprintf(os.Stderr, "\n")

	Server.Run()
}
