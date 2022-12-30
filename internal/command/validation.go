// SPDX-License-Identifier: Apache-2.0
// Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>
package command

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	_c "github.com/unrob/milpa/internal/constants"
)

type varSearchMap struct {
	Status int
	Name   string
	Usage  string
}

func (cmd *Command) Validate() (report map[string]int) {
	report = map[string]int{}

	for _, issue := range cmd.issues {
		report[issue.Error()] = 1
	}

	validate := validator.New()
	if err := validate.Struct(cmd); err != nil {
		verrs := err.(validator.ValidationErrors)
		for _, issue := range verrs {
			// todo: output better errors, see validator.FieldError
			report[fmt.Sprint(issue)] = 1
		}
	}

	if cmd.Meta.Kind != "source" {
		return report
	}

	vars := map[string]map[string]*varSearchMap{
		"argument": {},
		"option":   {},
	}

	for _, arg := range cmd.Arguments {
		vars["argument"][strings.ToUpper(strings.ReplaceAll(arg.Name, "-", "_"))] = &varSearchMap{2, arg.Name, ""}
	}

	for name := range cmd.Options {
		vars["option"][strings.ToUpper(strings.ReplaceAll(name, "-", "_"))] = &varSearchMap{2, name, ""}
	}

	contents, err := os.ReadFile(cmd.Meta.Path)
	if err != nil {
		report["Could not read script source"] = 1
		return
	}

	matches := _c.OutputPrefixPattern.FindAllStringSubmatch(string(contents), -1)
	for _, match := range matches {
		varName := match[len(match)-1]
		varKind := match[len(match)-2]

		kind := ""
		if varKind == "OPT" {
			kind = "option"
		} else if varKind == "ARG" {
			kind = "argument"
		}
		haystack := vars[kind]

		_, scriptVarIsValid := haystack[varName]
		if !scriptVarIsValid {
			haystack[varName] = &varSearchMap{Status: 1, Name: varName, Usage: match[0]}
		} else {
			haystack[varName].Status = 0
		}
	}

	for kind, col := range vars {
		for _, thisVar := range col {
			message := ""
			switch thisVar.Status {
			case 0:
				message = fmt.Sprintf("%s '%s' is used", kind, thisVar.Name)
			case 1:
				message = fmt.Sprintf("%s '%s' is not present in the spec, but used in the script as '%s'", kind, thisVar.Name, thisVar.Usage)
			case 2:
				message = fmt.Sprintf("%s '%s' is present in the spec, but not used by the script", kind, thisVar.Name)
			default:
				message = fmt.Sprintf("Unknown status %d for %s '%s'", thisVar.Status, kind, thisVar.Name)
			}

			report[message] = thisVar.Status
		}
	}

	return report
}
