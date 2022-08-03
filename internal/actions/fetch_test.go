package actions_test

import (
	"net/url"
	"os"
	"testing"

	"github.com/hashicorp/go-getter"
	"github.com/sirupsen/logrus"
)

func TestSomething(t *testing.T) {
	fqurl, err := getter.Detect("git@github.com:unRob/milpa.git", os.Getenv("PWD"), getter.Detectors)
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Infof("Detected uri: %s", fqurl)

	uri, err := url.Parse(fqurl)
	if err != nil {
		logrus.Fatal(err)
	}

	if uri.Scheme == "file" {
		logrus.Fatal("Refusing to copy local folder")
	}

	if uri.Opaque != "" && uri.Opaque[0] == ':' {
		logrus.Infof("Unwrapping uri: %s", uri.Opaque[1:])
		uri2, err := url.Parse(uri.Opaque[1:])
		if err != nil {
			logrus.Fatal(err)
		}

		uri = uri2
	}

	logrus.Infof("uri is: %s", uri)

	if uri.String() != "ssh://git@github.com/unRob/milpa.git" {
		t.Fatal("Weird URI")
		return
	}
}
