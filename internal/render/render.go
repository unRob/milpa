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
package render

import (
	"bytes"
	"os"

	"github.com/charmbracelet/glamour"
	_c "github.com/unrob/milpa/internal/constants"
	"github.com/unrob/milpa/internal/runtime"
	"golang.org/x/term"
)

func addBackticks(str []byte) []byte {
	return bytes.ReplaceAll(str, []byte("﹅"), []byte("`"))
}

func Markdown(content []byte, withColor bool) ([]byte, error) {
	content = addBackticks(content)

	if runtime.UnstyledHelpEnabled() {
		return content, nil
	}

	width, _, err := term.GetSize(0)
	if err != nil {
		return content, err
	}

	var styleFunc glamour.TermRendererOption

	if withColor {
		style := os.Getenv(_c.EnvVarHelpStyle)
		switch style {
		case "dark":
			styleFunc = glamour.WithStandardStyle("dark")
		case "light":
			styleFunc = glamour.WithStandardStyle("light")
		default:
			styleFunc = glamour.WithStandardStyle("auto")
		}
	} else {
		styleFunc = glamour.WithStandardStyle("notty")
	}

	renderer, err := glamour.NewTermRenderer(
		styleFunc,
		glamour.WithEmoji(),
		glamour.WithWordWrap(width),
	)

	if err != nil {
		return content, err
	}

	return renderer.RenderBytes(content)
}
