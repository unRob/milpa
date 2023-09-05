// SPDX-License-Identifier: Apache-2.0
// Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>
package docs

import (
	"os"

	"github.com/charmbracelet/glamour"
	"github.com/unrob/milpa/internal/constants"
)

func stringptr(str string) *string {
	return &str
}

func uintPointer(number uint) *uint {
	return &number
}

type StyleName string

const (
	StyleDark  StyleName = "dark"
	StyleLight StyleName = "light"
	StylePlain StyleName = "plain"
)

func init() {
	var zero uint
	glamour.NoTTYStyleConfig.Document.Margin = &zero
	glamour.NoTTYStyleConfig.Document.StylePrimitive.Color = nil
	glamour.DarkStyleConfig.Document.Margin = &zero
	glamour.DarkStyleConfig.Document.StylePrimitive.Color = nil
	glamour.LightStyleConfig.Document.Margin = &zero
	glamour.LightStyleConfig.Document.StylePrimitive.Color = nil
	glamour.DarkStyleConfig.List.Margin = uintPointer(2)
	glamour.LightStyleConfig.List.Margin = uintPointer(2)
	glamour.NoTTYStyleConfig.List.Margin = uintPointer(2)

	if os.Getenv(constants.EnvVarColorBitDepth) != "truecolor" {
		// Apple's Terminal.app does not support "true color", which is sad.
		glamour.DarkStyleConfig.H1.StylePrimitive.Color = stringptr("193")
		glamour.DarkStyleConfig.H1.StylePrimitive.BackgroundColor = stringptr("22")
		glamour.DarkStyleConfig.Heading.StylePrimitive.Color = stringptr("193")
		glamour.DarkStyleConfig.Code.StylePrimitive.Color = stringptr("230")
		glamour.DarkStyleConfig.Code.StylePrimitive.BackgroundColor = stringptr("22")

		glamour.LightStyleConfig.H1.StylePrimitive.Color = stringptr("193")
		glamour.LightStyleConfig.H1.StylePrimitive.BackgroundColor = stringptr("22")
		glamour.LightStyleConfig.Heading.StylePrimitive.Color = stringptr("28")
		glamour.LightStyleConfig.Code.StylePrimitive.Color = stringptr("22")
		glamour.LightStyleConfig.Code.StylePrimitive.BackgroundColor = stringptr("194")
		return
	}

	glamour.DarkStyleConfig.H1.StylePrimitive.Color = stringptr("#cefcd3")
	glamour.DarkStyleConfig.H1.StylePrimitive.BackgroundColor = stringptr("#2b3c2d")
	glamour.DarkStyleConfig.Heading.StylePrimitive.Color = stringptr("#c0e394")
	glamour.DarkStyleConfig.Code.StylePrimitive.Color = stringptr("#96b452")
	glamour.DarkStyleConfig.Code.StylePrimitive.BackgroundColor = stringptr("#132b17")

	glamour.LightStyleConfig.H1.StylePrimitive.Color = stringptr("#cefcd3")
	glamour.LightStyleConfig.H1.StylePrimitive.BackgroundColor = stringptr("#2b3c2d")
	glamour.LightStyleConfig.Heading.StylePrimitive.Color = stringptr("#12731D")
	glamour.LightStyleConfig.Code.StylePrimitive.Color = stringptr("#12731D")
	glamour.LightStyleConfig.Code.StylePrimitive.BackgroundColor = stringptr("#cee3c4")
}
