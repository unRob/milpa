// Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package runtime_test

import (
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"testing"

	_c "github.com/unrob/milpa/internal/constants"
	. "github.com/unrob/milpa/internal/runtime"
)

func TestCheckMilpaPathSet(t *testing.T) {
	MilpaPath = []string{"a", "b"}

	if err := CheckMilpaPathSet(); err != nil {
		t.Fatalf("Got error with set MILPA_PATH: %v", err)
	}

	MilpaPath = []string{}
	if err := CheckMilpaPathSet(); err == nil {
		t.Fatalf("Got no error with unset MILPA_PATH")
	}
}

func TestEnabled(t *testing.T) {
	defer func() { os.Setenv(_c.EnvVarMilpaVerbose, "") }()

	cases := []struct {
		Name    string
		Func    func() bool
		Expects bool
	}{
		{
			Name:    _c.EnvVarMilpaVerbose,
			Func:    VerboseEnabled,
			Expects: true,
		},
		{
			Name: _c.EnvVarValidationDisabled,
			Func: ValidationEnabled,
		},
		{
			Name: _c.EnvVarMilpaUnstyled,
			Func: ColorEnabled,
		},
		{
			Name: _c.EnvVarHelpUnstyled,
			Func: ColorEnabled,
		},
		{
			Name:    _c.EnvVarDebug,
			Func:    DebugEnabled,
			Expects: true,
		},
		{
			Name:    _c.EnvVarHelpUnstyled,
			Func:    UnstyledHelpEnabled,
			Expects: true,
		},
	}

	for _, c := range cases {
		fname := runtime.FuncForPC(reflect.ValueOf(c.Func).Pointer()).Name()
		name := fmt.Sprintf("%v/%s", fname, c.Name)
		enabled := []string{
			"yes", "true", "1", "enabled",
		}
		for _, val := range enabled {
			t.Run("enabled-"+val, func(t *testing.T) {
				os.Setenv(c.Name, val)
				if c.Func() != c.Expects {
					t.Fatalf("%s wasn't enabled with a valid value: %s", name, val)
				}
			})
		}

		disabled := []string{"", "no", "false", "0", "disabled"}
		for _, val := range disabled {
			t.Run("disabled-"+val, func(t *testing.T) {
				os.Setenv(c.Name, val)
				if c.Func() == c.Expects {
					t.Fatalf("%s was enabled with falsy value: %s", name, val)
				}
			})
		}
	}
}

func TestDoctorMode(t *testing.T) {
	cases := []struct {
		Args    []string
		Expects bool
	}{
		{
			Args: []string{},
		},
		{
			Args: []string{""},
		},
		{
			Args: []string{"something", "doctor"},
		},
		{
			Args:    []string{"__doctor"},
			Expects: true,
		},
		{
			Args:    []string{"__doctor", "whatever"},
			Expects: true,
		},
		{
			Args:    []string{"itself", "doctor"},
			Expects: true,
		},
	}

	for _, c := range cases {
		t.Run(strings.Join(c.Args, " "), func(t *testing.T) {
			os.Args = append([]string{"compa"}, c.Args...)
			res := DoctorModeEnabled()
			if res != c.Expects {
				t.Fatalf("Expected %v for %v and got %v", c.Expects, c.Args, res)
			}
		})
	}
}

func TestEnvironmentMapEnabled(t *testing.T) {
	MilpaPath = []string{"something"}
	trueString := strconv.FormatBool(true)
	os.Setenv(_c.EnvVarMilpaForceColor, trueString)
	os.Setenv(_c.EnvVarDebug, trueString)
	os.Setenv(_c.EnvVarMilpaVerbose, trueString)

	res := EnvironmentMap()
	if res == nil {
		t.Fatalf("Expected map, got nil")
	}

	expected := map[string]string{
		_c.EnvVarMilpaPath:       "something",
		_c.EnvVarMilpaForceColor: "always",
		_c.EnvVarMilpaPathParsed: trueString,
		_c.EnvVarDebug:           trueString,
		_c.EnvVarMilpaVerbose:    trueString,
	}

	if !reflect.DeepEqual(res, expected) {
		t.Fatalf("Unexpected result from enabled environment. Wanted %v, got %v", res, expected)
	}
}

func TestEnvironmentMapDisabled(t *testing.T) {
	MilpaPath = []string{"something"}
	trueString := strconv.FormatBool(true)
	// clear COLOR
	os.Unsetenv(_c.EnvVarMilpaForceColor)
	// set NO_COLOR
	os.Setenv(_c.EnvVarMilpaUnstyled, trueString)
	os.Unsetenv(_c.EnvVarDebug)
	os.Unsetenv(_c.EnvVarMilpaVerbose)
	os.Setenv(_c.EnvVarMilpaSilent, trueString)

	res := EnvironmentMap()
	if res == nil {
		t.Fatalf("Expected map, got nil")
	}

	expected := map[string]string{
		_c.EnvVarMilpaPath:       "something",
		_c.EnvVarMilpaUnstyled:   trueString,
		_c.EnvVarMilpaPathParsed: trueString,
		_c.EnvVarMilpaSilent:     trueString,
	}

	if !reflect.DeepEqual(res, expected) {
		t.Fatalf("Unexpected result from disabled environment. Wanted %v, got %v", res, expected)
	}
}
