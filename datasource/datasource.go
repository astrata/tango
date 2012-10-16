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

package datasource

import (
	"fmt"
	"github.com/astrata/tango/config"
	"github.com/gosexy/db"
	"github.com/gosexy/sugar"
	"math"
	"strconv"
)

// Number of items per page.
var ItemsPerPage = 9

// Creates a pager from a db.Collection.
func Pager(collection db.Collection, conds db.Cond, page int) sugar.Tuple {

	total, _ := collection.Count(conds)

	// Calculating total pages
	pages := int(math.Ceil(float64(total) / float64(ItemsPerPage)))

	// Getting current page
	if page < 1 {
		page = 1
	}

	prev := page - 1

	next := page + 1

	if next > total {
		next = 0
	}

	data := collection.FindAll(
		conds,
		db.Offset((page-1)*ItemsPerPage),
		db.Limit(ItemsPerPage),
	)

	response := sugar.Tuple{
		"total": total,
		"pager": sugar.Tuple{
			"count":   pages,
			"size":    ItemsPerPage,
			"next":    next,
			"prev":    prev,
			"current": page,
		},
		"data": data,
	}

	return response
}

// Returns a db.DataSource that you can use to connect to a database.
func Config(name string) (string, db.DataSource) {

	var driver string

	data := config.Get(fmt.Sprintf("datasource/%s", name))

	if data == nil {
		panic(fmt.Sprintf("tango: Cannot find %s datasource.", name))
	}

	source := db.DataSource{}

	for key, val := range data.(sugar.Tuple) {
		switch key {
		case "host":
			source.Host = fmt.Sprintf("%v", val)
		case "database":
			source.Database = fmt.Sprintf("%v", val)
		case "user":
			source.User = fmt.Sprintf("%v", val)
		case "password":
			source.Password = fmt.Sprintf("%v", val)
		case "port":
			port, _ := strconv.Atoi(fmt.Sprintf("%v", val))
			source.Port = port
		case "driver":
			driver = fmt.Sprintf("%v", val)
		}
	}

	return driver, source
}
