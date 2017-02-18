package iris2_test

import (
	"strconv"
	"testing"

	"github.com/go-iris2/iris2"
	"github.com/go-iris2/iris2/adaptors/gorillamux"
	"github.com/go-iris2/iris2/adaptors/httprouter"
	"github.com/go-iris2/iris2/httptest"
)

func testRouteStateSimple(t *testing.T, router iris2.Policy, offlineRoutePath string) {
	app := iris2.New()
	app.Adapt(router)

	offlineRouteRequestedTestPath := "/api/user/42"
	offlineBody := "user with id: 42"

	offlineRoute := app.None(offlineRoutePath, func(ctx *iris2.Context) {
		userid := ctx.Param("userid")
		if userid != "42" {
			// we are expecting userid 42 always in this test so
			t.Fatalf("what happened? expected userid to be 42 but got %s", userid)
		}
		ctx.Writef(offlineBody)
	}).ChangeName("api.users") // or an empty (), required, in order to get the Route instance.

	// change the "user.api" state from offline to online and online to offline
	app.Get("/change", func(ctx *iris2.Context) {
		// here
		if offlineRoute.IsOnline() {
			// set to offline
			app.Routes().Offline(offlineRoute)
		} else {
			// set to online if it was not online(so it was offline)
			app.Routes().Online(offlineRoute, iris2.MethodGet)
		}
	})

	app.Get("/execute", func(ctx *iris2.Context) {
		// here
		ctx.ExecRouteAgainst(offlineRoute, "/api/user/42")
	})

	// append the body and change the status code from an 'offline' route execution
	app.Get("/execute_modified", func(ctx *iris2.Context) {
		ctx.Set("mykey", "myval")
		// here
		ctx.Record() // if we want to control the response
		ctx.ExecRouteAgainst(offlineRoute, "/api/user/42")
		ctx.Write([]byte("modified from status code: " + strconv.Itoa(ctx.StatusCode())))
		ctx.SetStatusCode(iris2.StatusUseProxy)

		if ctx.Path() != "/execute_modified" {
			t.Fatalf("Expected Request Path of this context  NOT to change but got: '%s' ", ctx.Path())
		}

		if got := ctx.Get("mykey"); got != "myval" {
			t.Fatalf("Expected Value 'mykey' of this context  NOT to change('%s') but got: '%s' ", "myval", got)
		}
		ctx.Next()
	}, func(ctx *iris2.Context) {
		ctx.Writef("-original_middleware_here")
	})

	hello := "Hello from index"
	app.Get("/", func(ctx *iris2.Context) {
		ctx.Writef(hello)
	})

	e := httptest.New(app, t)

	e.GET("/").Expect().Status(iris2.StatusOK).Body().Equal(hello)
	// here
	// the status should be not found, the route is invisible from outside world
	e.GET(offlineRouteRequestedTestPath).Expect().Status(iris2.StatusNotFound)

	// set the route online with the /change
	e.GET("/change").Expect().Status(iris2.StatusOK)
	// try again, it should be online now
	e.GET(offlineRouteRequestedTestPath).Expect().Status(iris2.StatusOK).Body().Equal(offlineBody)
	// change to offline again
	e.GET("/change").Expect().Status(iris2.StatusOK)
	// and test again, it should be offline now
	e.GET(offlineRouteRequestedTestPath).Expect().Status(iris2.StatusNotFound)

	// finally test the execute on the offline route
	// it should be remains offline but execute the route like it is from client request.
	e.GET("/execute").Expect().Status(iris2.StatusOK).Body().Equal(offlineBody)
	e.GET(offlineRouteRequestedTestPath).Expect().Status(iris2.StatusNotFound)
	e.GET("/execute_modified").Expect().Status(iris2.StatusUseProxy).Body().
		Equal(offlineBody + "modified from status code: 200-original_middleware_here")
}

func TestRouteStateSimple(t *testing.T) {
	// httprouter adaptor
	testRouteStateSimple(t, httprouter.New(), "/api/user/:userid")
	// gorillamux adaptor
	testRouteStateSimple(t, gorillamux.New(), "/api/user/{userid}")
}
