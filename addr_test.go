// Black-box Testing
package iris2_test

import (
	"testing"

	"github.com/go-iris2/iris2"
)

func TestParseAddr(t *testing.T) {

	// test hosts
	expectedHost1 := "mydomain.com:1993"
	expectedHost2 := "mydomain.com"
	expectedHost3 := iris2.DefaultServerHostname + ":9090"
	expectedHost4 := "mydomain.com:443"

	host1 := iris2.ParseHost(expectedHost1)
	host2 := iris2.ParseHost(expectedHost2)
	host3 := iris2.ParseHost(":9090")
	host4 := iris2.ParseHost(expectedHost4)

	if host1 != expectedHost1 {
		t.Fatalf("Expecting server 1's host to be %s but we got %s", expectedHost1, host1)
	}
	if host2 != expectedHost2 {
		t.Fatalf("Expecting server 2's host to be %s but we got %s", expectedHost2, host2)
	}
	if host3 != expectedHost3 {
		t.Fatalf("Expecting server 3's host to be %s but we got %s", expectedHost3, host3)
	}
	if host4 != expectedHost4 {
		t.Fatalf("Expecting server 4's host to be %s but we got %s", expectedHost4, host4)
	}

	// test hostname
	expectedHostname1 := "mydomain.com"
	expectedHostname2 := "mydomain.com"
	expectedHostname3 := iris2.DefaultServerHostname
	expectedHostname4 := "mydomain.com"

	hostname1 := iris2.ParseHostname(host1)
	hostname2 := iris2.ParseHostname(host2)
	hostname3 := iris2.ParseHostname(host3)
	hostname4 := iris2.ParseHostname(host4)
	if hostname1 != expectedHostname1 {
		t.Fatalf("Expecting server 1's hostname to be %s but we got %s", expectedHostname1, hostname1)
	}

	if hostname2 != expectedHostname2 {
		t.Fatalf("Expecting server 2's hostname to be %s but we got %s", expectedHostname2, hostname2)
	}

	if hostname3 != expectedHostname3 {
		t.Fatalf("Expecting server 3's hostname to be %s but we got %s", expectedHostname3, hostname3)
	}

	if hostname4 != expectedHostname4 {
		t.Fatalf("Expecting server 4's hostname to be %s but we got %s", expectedHostname4, hostname4)
	}

	// test scheme, no need to test fullhost(scheme+host)
	expectedScheme1 := iris2.SchemeHTTP
	expectedScheme2 := iris2.SchemeHTTP
	expectedScheme3 := iris2.SchemeHTTP
	expectedScheme4 := iris2.SchemeHTTPS
	scheme1 := iris2.ParseScheme(host1)
	scheme2 := iris2.ParseScheme(host2)
	scheme3 := iris2.ParseScheme(host3)
	scheme4 := iris2.ParseScheme(host4)
	if scheme1 != expectedScheme1 {
		t.Fatalf("Expecting server 1's hostname to be %s but we got %s", expectedScheme1, scheme1)
	}

	if scheme2 != expectedScheme2 {
		t.Fatalf("Expecting server 2's hostname to be %s but we got %s", expectedScheme2, scheme2)
	}

	if scheme3 != expectedScheme3 {
		t.Fatalf("Expecting server 3's hostname to be %s but we got %s", expectedScheme3, scheme3)
	}

	if scheme4 != expectedScheme4 {
		t.Fatalf("Expecting server 4's hostname to be %s but we got %s", expectedScheme4, scheme4)
	}
}
