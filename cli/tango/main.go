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
	"flag"
	"fmt"
	"github.com/astrata/tango/version"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"reflect"
	"strings"
)

type Command struct {
	fn          interface{}
	Args        string
	Description string
}

var (
	PS = string(os.PathSeparator)

	goPath = ""

	tangoPath = ""
	tangoCwd  = ""

	flagHelp    = flag.Bool("help", false, "Shows help.")
	flagVersion = flag.Bool("version", false, "Shows version number.")

	commands = map[string]Command{}
)

func requireTango() {
	workdir := getCurrentWorkspace()

	if workdir == "" {
		die("Not a Tango! project (or any parent directory).")
	}
}

func goCommand(terms []string, targets []string) {
	var i int

	args := []string{}

	for i, _ = range terms {
		args = append(args, terms[i])
	}

	for i, _ = range targets {
		args = append(args, targets[i])
	}

	cmd := exec.Command("go", args...)

	cmd.Dir = getCurrentWorkspace()
	cmd.Stdin = os.Stdin

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	if err != nil {
		log.Printf("go: %s", err.Error())
	}
}

func die(message string) {
	fmt.Fprintln(os.Stderr, message)
	os.Exit(1)
}

func copyAll(src string, dst string) {

	info, err := os.Stat(src)

	if err != nil {
		panic(err)
	}

	if info.IsDir() == true {

		files, err := ioutil.ReadDir(src)

		if err != nil {
			panic(err)
		}

		for _, file := range files {
			info, _ = os.Stat(src + PS + file.Name())
			if info != nil {
				if info.IsDir() == true {
					os.Mkdir(dst+PS+file.Name(), os.FileMode(0775))
				}
				copyAll(src+PS+file.Name(), dst+PS+file.Name())
			}
		}

	} else {

		srcfp, err := os.Open(src)
		defer srcfp.Close()

		if err != nil {
			panic(err)
		}

		exists, _ := os.Stat(dst)

		if exists == nil {
			dstfp, err := os.Create(dst)
			defer dstfp.Close()

			if err != nil {
				panic(err)
			}

			io.Copy(dstfp, srcfp)
		}
	}

}

func checkUp() bool {
	tangoCwd, _ = os.Getwd()

	goPath = os.Getenv("GOPATH")

	paths := strings.Split(goPath, ":")

	if goPath == "" {
		log.Printf("Missing $GOPATH. A proper $GOPATH is required.")
		return false
	} else {

		for _, path := range paths {
			testPath := path + PS + "src" + PS + "github.com" + PS + "astrata" + PS + "tango"
			_, err := os.Stat(testPath + PS + "version" + PS + "version.go")
			if err == nil {
				tangoPath = testPath
				break
			}
		}

		if tangoPath == "" {
			log.Printf("Could not find tango installation directory.\n")
			return false
		}

	}

	tangoApp := getCurrentWorkspace()

	hijack := true

	if tangoApp != "" {
		for _, path := range paths {
			if path == tangoApp {
				hijack = false
			}
		}

		if hijack == true {
			goPath = goPath + ":" + tangoApp
			os.Setenv("GOPATH", goPath)
		}
	}

	return true
}

func banner() {
	fmt.Printf("Tango! (%s) - by Astrata\n", version.String)
	fmt.Println("")
}

func usage() {
	banner()

	fmt.Printf("Usage:\n\n")
	fmt.Printf("\ttango [command] [arguments]\n")

	fmt.Println("")

	fmt.Printf("The commands are:\n\n")

	for key, val := range commands {
		fmt.Printf("\t%-11s%s\n", key, val.Description)
	}
	fmt.Println("")
	fmt.Printf("Use \"tango help [command]\" for more information about a command.\n")
	fmt.Println("")

	//flag.PrintDefaults()
}

func main() {

	if checkUp() == false {
		return
	}

	flag.Parse()

	flag.Usage = usage

	args := flag.Args()

	argc := 0
	flags := true

	for _, arg := range args {
		if arg[0:1] == "-" && flags == true {

		} else {
			flags = false
			argc++
		}
	}

	if argc > 0 {

		command, exists := commands[args[0]]

		if exists {

			ctype := reflect.TypeOf(command.fn)

			eargs := ctype.NumIn()

			if (argc-1 == eargs) || (ctype.IsVariadic() && argc > 1) {

				vals := make([]reflect.Value, argc-1)

				for i := 1; i < argc; i++ {
					vals[i-1] = reflect.ValueOf(args[i])
				}

				fv := reflect.ValueOf(command.fn)
				fv.Call(vals)

				return

			} else {
				fmt.Printf("%s: Wrong number of arguments.\n", args[0])
			}

		} else {
			fmt.Printf("%s: Unknown command.\n", args[0])
		}

	}

	flag.Usage()
}
