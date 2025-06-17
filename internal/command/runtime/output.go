// SPDX-License-Identifier: Apache-2.0
// Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>
package runtime

import (
	"fmt"
	"strings"

	"github.com/alessio/shellescape"
	"github.com/unrob/milpa/internal/command/kind"
	"github.com/unrob/milpa/internal/command/meta"
	_c "github.com/unrob/milpa/internal/constants"
)

type valueSource interface {
	ToValue() any
	ToString() string
	Repeats() bool
}

func envVarValue(m meta.Meta, envName, value string) *string {
	switch m.Kind {
	case kind.Executable:
		value = fmt.Sprintf("%s=%s", envName, value)
	case kind.ShellScript:
		switch m.Shell {
		case kind.ShellBash, kind.ShellZSH:
			value = shellescape.Quote(value)
			value = fmt.Sprintf("export %s=%s", envName, value)
		case kind.ShellFish:
			value = fmt.Sprintf("set -x %s %s", envName, value)
		}
	}
	return &value
}

func envVarName(name string, prefix string) *string {
	if name == _c.HelpCommandName {
		return nil
	}
	if envName, ok := flagNames[name]; ok {
		return &envName
	}
	envName := fmt.Sprintf("%s%s", prefix, strings.ToUpper(strings.ReplaceAll(name, "-", "_")))
	return &envName
}

func envVarPair(name string, src valueSource, m meta.Meta) *string {
	if src.Repeats() {
		temp := []string{}
		var value string
		for _, v := range src.ToValue().([]string) {
			if m.Kind != kind.Executable {
				v = shellescape.Quote(v)
			}
			temp = append(temp, v)
		}
		switch m.Kind {
		case kind.Executable:
			value = fmt.Sprintf("%s=%s", name, strings.Join(temp, "|"))
		case kind.ShellScript:
			switch m.Shell {
			case kind.ShellBash, kind.ShellZSH:
				value = fmt.Sprintf("declare -a %s=( %s )", name, strings.Join(temp, " "))
			case kind.ShellFish:
				value = fmt.Sprintf("set -x %s %s", name, strings.Join(temp, " "))
			default:
				log.Fatalf("unhandled shell name: <%s>", m.Shell)
			}
		default:
			log.Fatalf("unhandled option kind: <%s>", m.Kind)
		}
		return &value
	}

	value := src.ToString()
	return envVarValue(m, name, value)
}
