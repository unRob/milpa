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
package render_test

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	_c "github.com/unrob/milpa/internal/constants"
	"github.com/unrob/milpa/internal/render"
)

func TestMarkdownUnstyled(t *testing.T) {
	content := []byte("# hello")
	os.Setenv(_c.EnvVarHelpUnstyled, "true")
	res, err := render.Markdown(content, false)

	if err != nil {
		t.Fatalf("Unexpected error %s", err)
	}

	expected := []byte("# hello") // nolint:ifshort
	if !reflect.DeepEqual(res, expected) {
		t.Fatalf("Unexpected response ---\n%s\n---\n wanted:\n---\n%s\n---", res, expected)
	}
}

func TestMarkdownNoColor(t *testing.T) {
	os.Unsetenv(_c.EnvVarHelpUnstyled)
	content := []byte("# hello ﹅world﹅")
	res, err := render.Markdown(content, false)

	if err != nil {
		t.Fatalf("Unexpected error %s", err)
	}

	// account for 80 character width word wrapping
	// our string is 15 characters, there's 2 for padding at the start
	spaces := "                                                             "

	expected := []byte("\n  # hello `world`" + spaces + "\n\n") // nolint:ifshort
	if !reflect.DeepEqual(res, expected) {
		t.Fatalf("Unexpected response ---\n%s\n---\n wanted:\n---\n%s\n---", res, expected)
	}
}

var autoStyleTestRender = "\n\x1b[38;5;228;48;5;63;1m\x1b[0m\x1b[38;5;228;48;5;63;1m\x1b[0m  \x1b[38;5;228;48;5;63;1m \x1b[0m\x1b[38;5;228;48;5;63;1mhello\x1b[0m\x1b[38;5;228;48;5;63;1m \x1b[0m\x1b[38;5;252m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[38;5;252m \x1b[0m\x1b[0m\n\x1b[0m\n"

const lightStyleTestRender = "\n\x1b[38;5;228;48;5;63;1m\x1b[0m\x1b[38;5;228;48;5;63;1m\x1b[0m  \x1b[38;5;228;48;5;63;1m \x1b[0m\x1b[38;5;228;48;5;63;1mhello\x1b[0m\x1b[38;5;228;48;5;63;1m \x1b[0m\x1b[38;5;234m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[38;5;234m \x1b[0m\x1b[0m\n\x1b[0m\n"

func TestMarkdownColor(t *testing.T) {
	os.Unsetenv(_c.EnvVarHelpUnstyled)
	content := []byte("# hello")

	styles := map[string][]byte{
		"":      []byte(autoStyleTestRender),
		"dark":  []byte(autoStyleTestRender),
		"auto":  []byte(autoStyleTestRender),
		"light": []byte(lightStyleTestRender),
	}
	for style, expected := range styles {
		t.Run(fmt.Sprintf("style %s", style), func(t *testing.T) {
			os.Setenv(_c.EnvVarHelpStyle, style)
			res, err := render.Markdown(content, true)

			if err != nil {
				t.Fatalf("Unexpected error %s", err)
			}

			if !reflect.DeepEqual(res, expected) {
				t.Fatalf("Unexpected response ---\n%v\n---\n wanted:\n---\n%v\n---", res, expected)
			}
		})
	}
}
