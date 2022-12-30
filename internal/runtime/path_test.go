// SPDX-License-Identifier: Apache-2.0
// Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>
package runtime_test

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	. "github.com/unrob/milpa/internal/runtime"
)

func testdataPathBuilder() func(string) string {
	wd, _ := os.Getwd()
	for !strings.HasSuffix(wd, "milpa") {
		wd = filepath.Dir(wd)
	}

	return func(suffix string) string {
		return fmt.Sprintf("%s/internal/runtime/testdata/%s", wd, suffix)
	}
}

func TestIsDir(t *testing.T) {
	tdp := testdataPathBuilder()

	if !IsDir(tdp("layer0/uno"), false) {
		t.Fatalf("Real directory is not a dir!")
	}

	buff := &bytes.Buffer{}
	fakePath := tdp("layer0/cuarenta-y-dos")
	logrus.SetOutput(buff)
	if IsDir(fakePath, true) {
		t.Fatalf("Fake directory marked as real")
	}

	warn := buff.String()
	text := fmt.Sprintf("Discarding non-directory <%s> from MILPA_PATH", fakePath)
	if !strings.Contains(warn, text) {
		t.Fatalf("Unexpected warning\n wanted %s\n got %s", text, warn)
	}

	if !strings.Contains(warn, "warning") {
		t.Fatalf("Unexpected warning\n wanted %s\n got %s", "warning", warn)
	}
}

func TestPathBuilderAddLookup(t *testing.T) {
	pb := NewPathBuilder()

	if pb == nil {
		t.Fatalf("Could not create pathbuilder")
	}

	startedWith := pb.LookupLen()
	someLookupFunc := func() []string {
		return []string{}
	}

	envVarName := "DISABLE_LOOKUP_FAKE"
	os.Setenv(envVarName, "true")
	pb.AddLookup(envVarName, someLookupFunc)

	got := pb.LookupLen()
	if got != startedWith {
		t.Fatalf("Found expected number of lookups: %d, found %d", startedWith, got)
	}

	startedWith = pb.LookupLen()
	os.Unsetenv(envVarName)
	pb.AddLookup(envVarName, someLookupFunc)
	got = pb.LookupLen()
	if got <= startedWith {
		t.Fatalf("Did not find expected number of lookups: %d, found %d", startedWith+1, got)
	}
}

func TestResolve(t *testing.T) {
	tdp := testdataPathBuilder()

	pb := &PathBuilder{}
	pb.AddLookup("primero", func() []string {
		return []string{tdp("layer0/uno"), tdp("layer0/dos"), tdp("layer0/tres"), tdp("layer2/uno-link")}
	})
	pb.AddLookup("second", func() []string {
		return []string{tdp("layer1/one"), tdp("layer1/two"), tdp("layer1/three")}
	})
	pb.AddLookup("third", func() []string {
		return []string{tdp("layer1/one"), tdp("layer1/two"), tdp("layer1/three")}
	})
	res := pb.Ordered()

	// they're in calling order, but each lookup is sorted alphabetically
	expected := []string{tdp("layer0/dos"), tdp("layer0/tres"), tdp("layer0/uno"), tdp("layer1/one"), tdp("layer1/three"), tdp("layer1/two")}
	if !reflect.DeepEqual(res, expected) {
		t.Fatalf("unexpected result, wanted: %v. got: %s", expected, res)
	}

	res = pb.Ordered()
	if !reflect.DeepEqual(res, expected) {
		t.Fatalf("unexpected result on second resolve, wanted: %v. got: %s", res, expected)
	}
}
