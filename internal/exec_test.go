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
package internal_test

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/spf13/cobra"
	. "github.com/unrob/milpa/internal"
)

func subshellSleep(ctx context.Context, env []string, executable string, args ...string) (bytes.Buffer, bytes.Buffer, error) {
	time.Sleep(100 * time.Nanosecond)
	return bytes.Buffer{}, bytes.Buffer{}, context.DeadlineExceeded
}

func subshellSucceed(ctx context.Context, env []string, executable string, args ...string) (bytes.Buffer, bytes.Buffer, error) {
	var out bytes.Buffer
	fmt.Fprint(&out, strings.Join([]string{
		"a",
		"b",
		"c",
	}, "\n"))
	return out, bytes.Buffer{}, nil
}

func TestExecTimesOut(t *testing.T) {
	ExecFunc = subshellSleep
	_, _, err := Exec("test-command", []string{"bash", "-c", "sleep", "2"}, 10*time.Nanosecond)
	if err == nil {
		t.Fatalf("timeout didn't happen after 10ms: %v", err)
	}
}

func TestExecWorksFine(t *testing.T) {
	ExecFunc = subshellSucceed
	args := []string{"a", "b", "c"}
	res, directive, err := Exec("test-command", append([]string{"bash", "-c", "echo"}, args...), 1*time.Second)
	if err != nil {
		t.Fatalf("good command failed: %v", err)
	}

	if directive != 0 {
		t.Fatalf("good command resulted in wrong directive, expected %d, got %d", 0, directive)
	}

	if strings.Join(args, "-") != strings.Join(res, "-") {
		t.Fatalf("good command resulted in wrong results, expected %v, got %v", res, args)
	}
}

func subshellError(ctx context.Context, env []string, executable string, args ...string) (bytes.Buffer, bytes.Buffer, error) {
	return bytes.Buffer{}, bytes.Buffer{}, fmt.Errorf("bad command is bad")
}

func TestExecErrors(t *testing.T) {
	ExecFunc = subshellError
	res, directive, err := Exec("test-command", []string{"bash", "-c", "bad-command"}, 1*time.Second)
	if err == fmt.Errorf("bad command is bad") {
		t.Fatalf("bad command didn't fail: %v", res)
	}

	if directive != cobra.ShellCompDirectiveError {
		t.Fatalf("bad command resulted in wrong directive, expected %d, got %d", cobra.ShellCompDirectiveError, directive)
	}

	if len(res) > 0 {
		t.Fatalf("bad command returned values, got %v", res)
	}
}
