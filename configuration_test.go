// Black-box Testing
package iris2_test

import (
	"reflect"
	"testing"

	. "github.com/go-iris2/iris2"
)

// $ go test -v -run TestConfiguration*

func TestConfigurationStatic(t *testing.T) {
	def := DefaultConfiguration()

	app := New(def)
	afterNew := *app.Config

	if !reflect.DeepEqual(def, afterNew) {
		t.Fatalf("Default configuration is not the same after NewFromConfig expected:\n %#v \ngot:\n %#v", def, afterNew)
	}

	afterNew.Charset = "changed"

	if reflect.DeepEqual(def, afterNew) {
		t.Fatalf("Configuration should be not equal, got: %#v", afterNew)
	}

	app = New(Configuration{DisableBodyConsumptionOnUnmarshal: true})

	afterNew = *app.Config

	if app.Config.DisableBodyConsumptionOnUnmarshal == false {
		t.Fatalf("Passing a Configuration field as Option fails, expected DisableBodyConsumptionOnUnmarshal to be true but was false")
	}

	app = New() // empty , means defaults so
	if !reflect.DeepEqual(def, *app.Config) {
		t.Fatalf("Default configuration is not the same after NewFromConfig expected:\n %#v \ngot:\n %#v", def, *app.Config)
	}
}

func TestConfigurationOptions(t *testing.T) {
	charset := "MYCHARSET"
	disableBanner := true

	app := New(OptionCharset(charset), OptionDisableBodyConsumptionOnUnmarshal(disableBanner))

	if got := app.Config.Charset; got != charset {
		t.Fatalf("Expected configuration Charset to be: %s but got: %s", charset, got)
	}

	if got := app.Config.DisableBodyConsumptionOnUnmarshal; got != disableBanner {
		t.Fatalf("Expected configuration DisableBodyConsumptionOnUnmarshal to be: %#v but got: %#v", disableBanner, got)
	}

	// now check if other default values are setted (should be setted automatically)

	expected := DefaultConfiguration()
	expected.Charset = charset
	expected.DisableBodyConsumptionOnUnmarshal = disableBanner

	has := *app.Config
	if !reflect.DeepEqual(has, expected) {
		t.Fatalf("Default configuration is not the same after New expected:\n %#v \ngot:\n %#v", expected, has)
	}
}

func TestConfigurationOptionsDeep(t *testing.T) {
	charset := "MYCHARSET"
	disableBanner := true
	vhost := "mydomain.com"
	// first charset,disableBanner and profilepath, no canonical order.
	app := New(OptionCharset(charset), OptionDisableBodyConsumptionOnUnmarshal(disableBanner), OptionVHost(vhost))

	expected := DefaultConfiguration()
	expected.Charset = charset
	expected.DisableBodyConsumptionOnUnmarshal = disableBanner
	expected.VHost = vhost

	has := *app.Config

	if !reflect.DeepEqual(has, expected) {
		t.Fatalf("DEEP configuration is not the same after New expected:\n %#v \ngot:\n %#v", expected, has)
	}
}
