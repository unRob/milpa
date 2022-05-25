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
package command

import (
	"fmt"
	"io/ioutil"
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
		report[issue] = 1
	}

	validate := validator.New()
	err := validate.Struct(cmd)
	if err != nil {
		verrs := err.(validator.ValidationErrors)
		for _, issue := range verrs {
			report[fmt.Sprint(issue)] = 1
		}
	}

	if cmd.Meta.Kind == "source" {
		contents, err := ioutil.ReadFile(cmd.Meta.Path)
		if err != nil {
			report["Could not read source"] = 1
			return
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
					message = fmt.Sprintf("%s '%s' is used but not defined, declared as '%s'", kind, thisVar.Name, thisVar.Usage)
				case 2:
					message = fmt.Sprintf("%s '%s' is not used but defined", kind, thisVar.Name)
				default:
					message = fmt.Sprintf("Unknown status %d for %s '%s'", thisVar.Status, kind, thisVar.Name)
				}

				report[message] = thisVar.Status
			}
		}
	}

	return report
}
