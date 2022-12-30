// SPDX-License-Identifier: Apache-2.0
// Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>
package render

import (
	"bytes"
	"os"

	"github.com/charmbracelet/glamour"
	"github.com/sirupsen/logrus"
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
		logrus.Debugf("Could not get terminal width")
		width = 80
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
