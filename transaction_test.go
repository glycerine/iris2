package iris2_test

import (
	"github.com/go-iris2/iris2"
	"github.com/go-iris2/iris2/httptest"

	"net/http"
	"testing"
)

func TestTransaction(t *testing.T) {
	app := iris2.New()
	app.Adapt(newTestNativeRouter())

	firstTransactionFailureMessage := "Error: Virtual failure!!!"
	secondTransactionSuccessHTMLMessage := "<h1>This will sent at all cases because it lives on different transaction and it doesn't fails</h1>"
	persistMessage := "<h1>I persist show this message to the client!</h1>"

	maybeFailureTransaction := func(shouldFail bool, isRequestScoped bool) func(t *iris2.Transaction) {
		return func(t *iris2.Transaction) {
			// OPTIONAl, the next transactions and the flow will not be skipped if this transaction fails
			if isRequestScoped {
				t.SetScope(iris2.RequestTransactionScope)
			}

			// OPTIONAL STEP:
			// create a new custom type of error here to keep track of the status code and reason message
			err := iris2.NewTransactionErrResult()

			t.Context.Text(http.StatusOK, "Blablabla this should not be sent to the client because we will fill the err with a message and status")

			fail := shouldFail

			if fail {
				err.StatusCode = http.StatusInternalServerError
				err.Reason = firstTransactionFailureMessage
			}

			// OPTIONAl STEP:
			// but useful if we want to post back an error message to the client if the transaction failed.
			// if the reason is empty then the transaction completed successfully,
			// otherwise we rollback the whole response body and cookies and everything lives inside the transaction.Request.
			t.Complete(err)
		}
	}

	successTransaction := func(scope *iris2.Transaction) {
		if scope.Context.Request.RequestURI == "/failAllBecauseOfRequestScopeAndFailure" {
			t.Fatalf("We are inside successTransaction but the previous REQUEST SCOPED TRANSACTION HAS FAILED SO THiS SHOULD NOT BE RAN AT ALL")
		}
		scope.Context.HTML(http.StatusOK,
			secondTransactionSuccessHTMLMessage)
		// * if we don't have any 'throw error' logic then no need of scope.Complete()
	}

	persistMessageHandler := func(ctx *iris2.Context) {
		// OPTIONAL, depends on the usage:
		// at any case, what ever happens inside the context's transactions send this to the client
		ctx.HTML(http.StatusOK, persistMessage)
	}

	app.Get("/failFirsTransactionButSuccessSecondWithPersistMessage", func(ctx *iris2.Context) {
		ctx.BeginTransaction(maybeFailureTransaction(true, false))
		ctx.BeginTransaction(successTransaction)
		persistMessageHandler(ctx)
	})

	app.Get("/failFirsTransactionButSuccessSecond", func(ctx *iris2.Context) {
		ctx.BeginTransaction(maybeFailureTransaction(true, false))
		ctx.BeginTransaction(successTransaction)
	})

	app.Get("/failAllBecauseOfRequestScopeAndFailure", func(ctx *iris2.Context) {
		ctx.BeginTransaction(maybeFailureTransaction(true, true))
		ctx.BeginTransaction(successTransaction)
	})

	customErrorTemplateText := "<h1>custom error</h1>"
	app.OnError(http.StatusInternalServerError, func(ctx *iris2.Context) {
		ctx.Text(http.StatusInternalServerError, customErrorTemplateText)
	})

	failureWithRegisteredErrorHandler := func(ctx *iris2.Context) {
		ctx.BeginTransaction(func(transaction *iris2.Transaction) {
			transaction.SetScope(iris2.RequestTransactionScope)
			err := iris2.NewTransactionErrResult()
			err.StatusCode = http.StatusInternalServerError // set only the status code in order to execute the registered template
			transaction.Complete(err)
		})

		ctx.Text(http.StatusOK, "this will not be sent to the client because first is requested scope and it's failed")
	}

	app.Get("/failAllBecauseFirstTransactionFailedWithRegisteredErrorTemplate", failureWithRegisteredErrorHandler)

	e := httptest.New(app, t)

	e.GET("/failFirsTransactionButSuccessSecondWithPersistMessage").
		Expect().
		Status(http.StatusOK).
		ContentType("text/html", app.Config.Charset).
		Body().
		Equal(secondTransactionSuccessHTMLMessage + persistMessage)

	e.GET("/failFirsTransactionButSuccessSecond").
		Expect().
		Status(http.StatusOK).
		ContentType("text/html", app.Config.Charset).
		Body().
		Equal(secondTransactionSuccessHTMLMessage)

	e.GET("/failAllBecauseOfRequestScopeAndFailure").
		Expect().
		Status(http.StatusInternalServerError).
		Body().
		Equal(firstTransactionFailureMessage)

	e.GET("/failAllBecauseFirstTransactionFailedWithRegisteredErrorTemplate").
		Expect().
		Status(http.StatusInternalServerError).
		Body().
		Equal(customErrorTemplateText)
}
