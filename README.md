Iris2
=====

[![Build Status](https://travis-ci.org/go-iris2/iris2.svg?branch=master)](https://travis-ci.org/go-iris2/iris2)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-iris2/iris2)](https://goreportcard.com/report/github.com/go-iris2/iris2)
[![GoDoc](https://godoc.org/github.com/go-iris2/iris2?status.svg)](https://godoc.org/github.com/go-iris2/iris2)
[![codecov](https://codecov.io/gh/go-iris2/iris2/branch/master/graph/badge.svg)](https://codecov.io/gh/go-iris2/iris2)


Iris2 is a fork of the original [Iris](https://github.com/kataras/iris) framework. As Iris, Iris2 is an efficient and complete toolbox with robust set of features.<br/>Write <b>your own</b>
<b>perfect high-performance web applications</b> <br/>with unlimited potentials and <b>portability</b>.<br/>
Iris2 aims for a *stable* API and be easy to use.

Installation
------------

The only requirement is the [Go Programming Language](https://golang.org/dl/), version 1.8+

```bash
$ go get github.com/go-iris2/iris2
```


Overview
--------

```go
package main

import (
	"github.com/go-iris2/iris2"
	"github.com/go-iris2/iris2/adaptors/view"
)

func main() {
	app := iris.New()
	app.Adapt(iris.Devlogger())

	// 5 template engines are supported out-of-the-box:
	//
	// - standard html/template
	// - amber
	// - django
	// - handlebars
	// - pug(jade)
	//
	// Use the html standard engine for all files inside "./views" folder with extension ".html"
	templates := view.HTML("./views", ".html")
	app.Adapt(templates)

	// http://localhost:6200
	// Method: "GET"
	// Render ./views/index.html
	app.Get("/", func(ctx *iris.Context) {
		ctx.Render("index.html", nil)
	})

	// Group routes, optionally: share middleware, template layout and custom http errors.
	userAPI := app.Party("/users", userAPIMiddleware).
		Layout("layouts/userLayout.html")
	{
		// Fire userNotFoundHandler when Not Found
		// inside http://localhost:6200/users/*anything
		userAPI.OnError(404, userNotFoundHandler)

		// http://localhost:6200/users
		// Method: "GET"
		userAPI.Get("/", getAllHandler)

		// http://localhost:6200/users/42
		// Method: "GET"
		userAPI.Get("/:id", getByIDHandler)

		// http://localhost:6200/users
		// Method: "POST"
		userAPI.Post("/", saveUserHandler)
	}

	// Start the server at 127.0.0.1:6200
	app.Listen(":6200")
}

func getByIDHandler(ctx *iris.Context) {
	// take the :id from the path, parse to integer
	// and set it to the new userID local variable.
	userID, _ := ctx.ParamInt("id")

	// userRepo, imaginary database service <- your only job.
	user := userRepo.GetByID(userID)

	// send back a response to the client,
	// .JSON: content type as application/json; charset="utf-8"
	// iris.StatusOK: with 200 http status code.
	//
	// send user as it is or make use of any json valid golang type,
	// like the iris.Map{"username" : user.Username}.
	ctx.JSON(iris.StatusOK, user)
}
```

Versioning
----------

Iris2 adheres to [Semantic Versioning](http://semver.org/).

See the [Changelog](CHANGELOG.md).

License
-------

Unless otherwise noted, the source files are distributed
under the MIT License found in the [LICENSE](LICENSE).
