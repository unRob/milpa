// SPDX-License-Identifier: Apache-2.0
// Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>
package docs

import (
	"os"

	"github.com/charmbracelet/glamour/ansi"
	"github.com/charmbracelet/glamour/styles"

	"github.com/unrob/milpa/internal/constants"
)

func stringptr(str string) *string { return &str }

func uintPointer(number uint) *uint { return &number }

func boolptr(b bool) *bool { return &b }

var MilpaNoTTY = ansi.StyleConfig{
	Document: ansi.StyleBlock{
		Margin: uintPointer(0),
		StylePrimitive: ansi.StylePrimitive{
			Color:       nil,
			BlockPrefix: "\n",
			BlockSuffix: "\n",
		},
	},
	List: ansi.StyleList{
		LevelIndent: 2,
		StyleBlock: ansi.StyleBlock{
			Margin: uintPointer(2),
		},
	},
	Item: ansi.StylePrimitive{
		BlockPrefix: "• ",
	},
	Strikethrough: ansi.StylePrimitive{
		CrossedOut: boolptr(true),
	},
	Emph: ansi.StylePrimitive{
		Italic: boolptr(false),
	},
	Strong: ansi.StylePrimitive{
		Bold: boolptr(false),
	},
	Link: ansi.StylePrimitive{
		Underline: boolptr(false),
	},
	LinkText: ansi.StylePrimitive{
		Bold: boolptr(false),
	},
	H1: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Bold:   boolptr(false),
			Prefix: "# ",
			Suffix: "",
		},
	},
	H2: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Prefix: "## ",
		},
	},
	H3: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Prefix: "### ",
		},
	},
	H4: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Prefix: "#### ",
		},
	},
	H5: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Prefix: "##### ",
		},
	},
	H6: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Prefix: "###### ",
			Bold:   boolptr(false),
		},
	},
	Heading: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			BlockSuffix: "\n",
			Bold:        boolptr(false),
		},
	},
	Code: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Prefix: "`",
			Suffix: "`",
		},
	},
}

var MilpaTrueColorDark = ansi.StyleConfig{
	Document: ansi.StyleBlock{
		Margin: uintPointer(0),
		StylePrimitive: ansi.StylePrimitive{
			Color:       nil,
			BlockPrefix: "\n",
			BlockSuffix: "\n",
		},
	},
	List: ansi.StyleList{
		LevelIndent: 2,
		StyleBlock: ansi.StyleBlock{
			Margin: uintPointer(2),
		},
	},
	Item: ansi.StylePrimitive{
		BlockPrefix: "• ",
	},
	Strikethrough: ansi.StylePrimitive{
		CrossedOut: boolptr(true),
	},
	Emph: ansi.StylePrimitive{
		Italic: boolptr(true),
	},
	Strong: ansi.StylePrimitive{
		Bold: boolptr(true),
	},
	Link: ansi.StylePrimitive{
		Color:     stringptr("30"),
		Underline: boolptr(true),
	},
	LinkText: ansi.StylePrimitive{
		Color: stringptr("35"),
		Bold:  boolptr(true),
	},
	H1: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Color:           stringptr("#cefcd3"),
			BackgroundColor: stringptr("#2b3c2d"),
			Bold:            boolptr(true),
			Prefix:          " ",
			Suffix:          " ",
		},
	},
	H2: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Prefix: "## ",
		},
	},
	H3: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Prefix: "### ",
		},
	},
	H4: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Prefix: "#### ",
		},
	},
	H5: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Prefix: "##### ",
		},
	},
	H6: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Prefix: "###### ",
			Color:  stringptr("35"),
			Bold:   boolptr(false),
		},
	},
	Heading: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			BlockSuffix: "\n",
			Color:       stringptr("#c0e394"),
			Bold:        boolptr(true),
		},
	},
	Code: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Prefix:          " ",
			Suffix:          " ",
			Color:           stringptr("#96b452"),
			BackgroundColor: stringptr("#132b17"),
		},
	},
}

var MilpaTrueColorLight = ansi.StyleConfig{
	Document: ansi.StyleBlock{
		Margin: uintPointer(0),
		StylePrimitive: ansi.StylePrimitive{
			Color:       nil,
			BlockPrefix: "\n",
			BlockSuffix: "\n",
		},
	},
	List: ansi.StyleList{
		LevelIndent: 2,
		StyleBlock: ansi.StyleBlock{
			Margin: uintPointer(2),
		},
	},
	Item: ansi.StylePrimitive{
		BlockPrefix: "• ",
	},
	Strikethrough: ansi.StylePrimitive{
		CrossedOut: boolptr(true),
	},
	Emph: ansi.StylePrimitive{
		Italic: boolptr(true),
	},
	Strong: ansi.StylePrimitive{
		Bold: boolptr(true),
	},
	Link: ansi.StylePrimitive{
		Color:     stringptr("30"),
		Underline: boolptr(true),
	},
	LinkText: ansi.StylePrimitive{
		Color: stringptr("35"),
		Bold:  boolptr(true),
	},
	H1: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Color:           stringptr("#cefcd3"),
			BackgroundColor: stringptr("#2b3c2d"),
			Bold:            boolptr(true),
			Prefix:          " ",
			Suffix:          " ",
		},
	},
	H2: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Prefix: "## ",
		},
	},
	H3: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Prefix: "### ",
		},
	},
	H4: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Prefix: "#### ",
		},
	},
	H5: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Prefix: "##### ",
		},
	},
	H6: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Prefix: "###### ",
			Color:  stringptr("35"),
			Bold:   boolptr(false),
		},
	},
	Heading: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			BlockSuffix: "\n",
			Color:       stringptr("#12731D"),
			Bold:        boolptr(true),
		},
	},
	Code: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Prefix:          " ",
			Suffix:          " ",
			Color:           stringptr("#12731D"),
			BackgroundColor: stringptr("#cee3c4"),
		},
	},
}

var MilpaDark = ansi.StyleConfig{
	Document: ansi.StyleBlock{
		Margin: uintPointer(0),
		StylePrimitive: ansi.StylePrimitive{
			Color:       nil,
			BlockPrefix: "\n",
			BlockSuffix: "\n",
		},
	},
	List: ansi.StyleList{
		StyleBlock: ansi.StyleBlock{
			Margin: uintPointer(2),
		},
	},
	H1: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Color:           stringptr("193"),
			BackgroundColor: stringptr("22"),
			Bold:            boolptr(true),
		},
	},
	Heading: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Color: stringptr("193"),
			Bold:  boolptr(true),
		},
	},
	Code: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Color:           stringptr("230"),
			BackgroundColor: stringptr("22"),
		},
	},
}

var MilpaLight = ansi.StyleConfig{
	Document: ansi.StyleBlock{
		Margin: uintPointer(0),
		StylePrimitive: ansi.StylePrimitive{
			Color:       nil,
			BlockPrefix: "\n",
			BlockSuffix: "\n",
		},
	},
	List: ansi.StyleList{
		StyleBlock: ansi.StyleBlock{
			Margin: uintPointer(2),
		},
	},
	H1: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Color:           stringptr("193"),
			BackgroundColor: stringptr("22"),
			Bold:            boolptr(true),
			Prefix:          " ",
			Suffix:          " ",
		},
	},
	H2: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Prefix: "## ",
		},
	},
	H3: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Prefix: "### ",
		},
	},
	H4: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Prefix: "#### ",
		},
	},
	H5: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Prefix: "##### ",
		},
	},
	H6: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Prefix: "###### ",
			Color:  stringptr("35"),
			Bold:   boolptr(false),
		},
	},
	Heading: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Color:       stringptr("28"),
			Bold:        boolptr(true),
			BlockSuffix: "\n",
		},
	},
	Code: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Color:           stringptr("22"),
			BackgroundColor: stringptr("194"),
		},
	},
}

func init() {
	styles.LightStyleConfig = MilpaTrueColorLight
	styles.DarkStyleConfig = MilpaTrueColorDark
	styles.NoTTYStyleConfig = MilpaNoTTY
	if os.Getenv(constants.EnvVarColorBitDepth) != "truecolor" {
		styles.LightStyleConfig = MilpaLight
		styles.DarkStyleConfig = MilpaDark
	}
}
