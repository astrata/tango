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

package main

import (
	"io/ioutil"
	"runtime"
	"strconv"
	"strings"
)

func init() {
	commands["build"] = Command{buildCommand, "<output>", "Compiles a Tango! app into a static binary."}
}

func getProjectName() string {
	workdir := getCurrentWorkspace()
	if workdir != "" {
		chunks := strings.Split(workdir, PS)
		return chunks[len(chunks)-1:][0]
	}
	return ""
}

func getGoFiles() []string {
	var apps []string

	workdir := getCurrentWorkspace()

	if workdir == "" {
		die("Not a Tango! app (or any parent directory).")
	}

	files, err := ioutil.ReadDir(workdir)

	if err != nil {
		die("Could not read working directory.")
	}

	for _, file := range files {
		if file.IsDir() == false {
			name := file.Name()
			if strings.HasSuffix(name, ".go") == true {
				if strings.HasSuffix(name, "_test.go") == false {
					apps = append(apps, workdir+PS+name)
				}
			}
		}
	}

	return apps
}

func buildCommand() {
	requireTango()

	apps := getGoFiles()

	name := getProjectName()

	goCommand([]string{"build", "-o", name, "-p", strconv.Itoa(runtime.NumCPU())}, apps)
}
