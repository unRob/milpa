// SPDX-License-Identifier: Apache-2.0
// Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>
package actions

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"

	"git.rob.mx/nidito/chinampa/pkg/command"
	"git.rob.mx/nidito/chinampa/pkg/env"
	"git.rob.mx/nidito/chinampa/pkg/errors"
	"git.rob.mx/nidito/chinampa/pkg/render"
	"git.rob.mx/nidito/chinampa/pkg/statuscode"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/unrob/milpa/internal/docs"

	milpaCommand "github.com/unrob/milpa/internal/command"
	_c "github.com/unrob/milpa/internal/constants"
	"github.com/unrob/milpa/internal/logger"
	"github.com/unrob/milpa/internal/lookup"
)

var dlog = logger.Sub("action:docs")

var AfterHelp = os.Exit

func startServer(listen, address string) error {
	logrus.Warnf("Using static resources at %s", os.Getenv("MILPA_ROOT")+"/internal/static")
	os.Setenv(env.HelpUnstyled, "true")
	http.Handle("/static/", docs.StaticHandler())
	http.HandleFunc("/", docs.RenderHandler(address))

	dlog.Info("Starting help http server")
	return http.ListenAndServe(listen, nil)
}

var Docs = &command.Command{
	Path:    []string{"help", "docs"},
	Summary: "Displays docs on TOPIC",
	Description: "Shows markdown-formatted documentation from milpa repos. See `" + _c.Milpa + " " + _c.HelpCommandName + " docs milpa repo docs` for more information on how to write your own." + `

An HTTP server to browse documentation can be started by running:

﹅﹅﹅sh
milpa help docs --server
﹅﹅﹅
`,
	Arguments: command.Arguments{
		&command.Argument{
			Name:        "topic",
			Description: "The topic to show docs for",
			Variadic:    true,
			Required:    false,
			Values: &command.ValueSource{
				Suggestion: true,
				Func: func(cmd *command.Command, currentValue, config string) (values []string, flag cobra.ShellCompDirective, err error) {
					args := cmd.Arguments[0].ToValue().([]string)
					dlog.Debugf("looking for docs given %v and %s", args, currentValue)

					cv := ""
					if len(args) > 1 {
						cv = args[len(args)-1]
						args = args[0 : len(args)-1]
					}
					dlog.Debugf("looking for docs given %v and %s", args, cv)
					docs, err := lookup.Docs(args, cv, false)
					if err != nil {
						return nil, cobra.ShellCompDirectiveNoFileComp, err
					}

					return docs, cobra.ShellCompDirectiveNoFileComp, nil
				},
			},
		},
	},
	Options: command.Options{
		"server": {
			Description: "Starts an http server at the specified address",
			Type:        command.ValueTypeBoolean,
			Default:     false,
		},
		"listen": {
			Description: "The address to listen at when using `--server`",
			Type:        command.ValueTypeString,
			Default:     "localhost:4242",
		},
		"base": {
			Description: "A URL base to use for rendering html links",
			Type:        command.ValueTypeString,
			Default:     "http://localhost:4242",
		},
	},
	Meta: milpaCommand.Meta{
		Path: os.Getenv(_c.EnvVarMilpaRoot) + "/milpa/docs",
		Name: []string{_c.HelpCommandName, "docs"},
		Repo: os.Getenv(_c.EnvVarMilpaRoot),
		Kind: "docs",
	},
	HelpFunc: func(printLinks bool) string {
		dlog.Debug("showing docs help")
		topics, err := lookup.Docs([]string{}, "", false)
		if err != nil {
			return ""
		}
		topicList := []string{}
		for _, topic := range topics {
			if printLinks {
				topic = fmt.Sprintf("[%s](%s)", topic, topic)
			}
			topicList = append(topicList, "- "+topic)
		}

		return "## Available topics:\n\n" + strings.Join(topicList, "\n")
	},
	Action: func(cmd *command.Command) error {
		dlog.Debug("Rendering docs")
		args := cmd.Arguments[0].ToValue().([]string)
		if len(args) == 0 {
			if cmd.Options["server"].ToValue().(bool) {
				listen := cmd.Options["listen"].ToString()
				address := cmd.Options["base"].ToString()
				dlog.Infof("Starting docs server at %s", listen)
				return startServer(listen, address)
			}
			dlog.Debug("Rendering docs help page")
			err := cmd.Cobra.Help()
			if err != nil {
				return err
			}
			AfterHelp(statuscode.RenderHelp)
			dlog.Debug("Rendered docs help page")
			return nil
		}

		contents, err := docs.FromQuery(args)
		if err != nil {
			switch err.(type) {
			case errors.BadArguments:
				return err
			}
			return errors.NotFound{Msg: "Unknown doc: " + err.Error()}
		}

		titleExp := regexp.MustCompile("^title: (.+)")
		frontmatterSep := []byte("---\n")
		if len(contents) > 3 && string(contents[0:4]) == string(frontmatterSep) {
			// strip out frontmatter
			parts := bytes.SplitN(contents, frontmatterSep, 3)
			title := titleExp.FindString(string(parts[1]))
			if title != "" {
				title = strings.TrimPrefix(title, "title: ")
			} else {
				title = strings.Join(args, " ")
			}
			contents = bytes.Join([][]byte{[]byte("# " + title + "\n"), parts[2]}, []byte("\n"))
		}

		withColor, _ := cmd.Cobra.Flags().GetBool("no-color")

		doc, err := render.Markdown(contents, !withColor)
		if err != nil {
			return err
		}

		if _, err := cmd.Cobra.OutOrStderr().Write(doc); err != nil {
			return err
		}
		AfterHelp(statuscode.RenderHelp)

		return nil
	},
}
