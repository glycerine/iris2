package errors

import (
	"fmt"
	"testing"
)

var errMessage = "User with mail: %s already exists"
var errUserAlreadyExists = New(errMessage)
var userMail = "user1@mail.go"
var expectedUserAlreadyExists = "User with mail: user1@mail.go already exists"

func ExampleError() {
	fmt.Print(errUserAlreadyExists.Format(userMail))
	fmt.Printf("\n")
	// first output first Output line
	fmt.Print(errUserAlreadyExists.Format(userMail).Append("\nPlease change your mail addr"))
	// second output second and third Output lines

	// Output:
	// User with mail: user1@mail.go already exists
	// User with mail: user1@mail.go already exists
	// Please change your mail addr
}

func do(method string, testErr *Error, expectingMsg string, t *testing.T) {
	formattedErr := func() error {
		return testErr.Format(userMail)
	}()

	if formattedErr.Error() != expectingMsg {
		t.Fatalf("Error %s failed, expected:\n%s got:\n%s", method, expectingMsg, formattedErr.Error())
	}
}

func TestFormat(t *testing.T) {
	expected := expectedUserAlreadyExists
	do("Format Test", errUserAlreadyExists, expected, t)
}

func TestAppendErr(t *testing.T) {
	errChangeMailMsg := "Please change your mail addr"
	errChangeMail := fmt.Errorf(errChangeMailMsg)                                                           // test go standard error
	expectedErrorMessage := errUserAlreadyExists.Format(userMail).Error() + errChangeMailMsg // first Prefix and last newline lives inside do
	errAppended := errUserAlreadyExists.AppendErr(errChangeMail)
	do("Append Test Standard error type", &errAppended, expectedErrorMessage, t)
}

func TestAppendError(t *testing.T) {
	errChangeMailMsg := "Please change your mail addr"
	errChangeMail := New(errChangeMailMsg)                                                                       // test Error struct
	expectedErrorMessage := errUserAlreadyExists.Format(userMail).Error() + errChangeMail.Error() // first Prefix and last newline lives inside do
	errAppended := errUserAlreadyExists.AppendErr(errChangeMail)
	do("Append Test Error type", &errAppended, expectedErrorMessage, t)
}

func TestAppend(t *testing.T) {
	errChangeMailMsg := "Please change your mail addr"
	expectedErrorMessage := errUserAlreadyExists.Format(userMail).Error() + errChangeMailMsg // first Prefix and last newline lives inside do
	errAppended := errUserAlreadyExists.Append(errChangeMailMsg)
	do("Append Test string Message", &errAppended, expectedErrorMessage, t)
}
