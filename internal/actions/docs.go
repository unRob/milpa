// SPDX-License-Identifier: Apache-2.0
// Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>
package actions

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"time"

	"git.rob.mx/nidito/chinampa/pkg/command"
	"git.rob.mx/nidito/chinampa/pkg/env"
	"git.rob.mx/nidito/chinampa/pkg/errors"
	"git.rob.mx/nidito/chinampa/pkg/render"
	"git.rob.mx/nidito/chinampa/pkg/statuscode"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/unrob/milpa/internal/docs"

	"git.rob.mx/nidito/chinampa/pkg/logger"
	milpaCommand "github.com/unrob/milpa/internal/command"
	_c "github.com/unrob/milpa/internal/constants"
	"github.com/unrob/milpa/internal/lookup"
)

var dlog = logger.Sub("action:docs")

var AfterHelp = os.Exit

func startServer(listen, address string) error {
	os.Setenv(env.HelpStyle, "markdown")
	// Replace with DevelopmentStaticResourceHandler to use locally available
	// static resources during development
	http.Handle("/static/", docs.EmbeddedStaticResourceHandler())
	http.HandleFunc("/", docs.RenderHandler(address))

	server := &http.Server{
		Addr:              listen,
		ReadHeaderTimeout: 3 * time.Second,
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		os.Exit(0)
	}()
	return server.ListenAndServe()
}

func init() {
	// This convoluted piece of code will ensure `help docs` render
	// the correct address in its own help
	// as well as rendering available docs topics
	Docs.HelpFunc = func(printLinks bool) string {
		base := Docs.Options["base"].ToString()
		listen := Docs.Options["listen"]
		defaultListen := listen.Default.(string)
		base = strings.ReplaceAll(base, defaultListen, listen.ToString())
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

		return `### Server mode

An HTTP server to browse documentation can be started by running:

﹅﹅﹅sh
milpa help docs --server
# then head to http://localhost:4242
﹅﹅﹅

Command and docs are available at their names, replacing spaces with forward slashes ﹅/﹅, for example:

- ` + base + `/help/docs will show this help page.
- ` + base + `/itself/doctor shows documentation for the ﹅milpa itself doctor﹅ command.
- ` + base + `/help/docs/milpa renders the file at ﹅.milpa/docs/milpa.md﹅ (or ﹅.milpa/docs/milpa/index.md﹅).

## Available topics

` + strings.Join(topicList, "\n")
	}
}

var Docs = &command.Command{
	Path:        []string{"help", "docs"},
	Summary:     "Displays docs on TOPIC",
	Description: "Shows markdown-formatted documentation from milpa repos. See `" + _c.Milpa + " " + _c.HelpCommandName + " docs milpa repo docs` for more information on how to write your own.",
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
					if len(args) > 1 && args[len(args)-1] == "" {
						// remove last argument from docs path lookup base
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
	Action: func(cmd *command.Command) error {
		args := cmd.Arguments[0].ToValue().([]string)
		if len(args) == 0 {
			if cmd.Options["server"].ToValue().(bool) {
				listen := cmd.Options["listen"].ToString()
				base := cmd.Options["base"]
				defaultListen := cmd.Options["listen"].Default.(string)
				address := strings.ReplaceAll(base.ToString(), defaultListen, listen)
				dlog.Infof("Starting docs server at http://%s, press CTRL-C to stop...", listen)
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
				// when a doc is not found, let's show docs help,
				// as error bubbling makes cobra think the calling command
				// is milpa (since "docs" is a child of "help")
				helpErr := cmd.Cobra.Help()
				if helpErr != nil {
					os.Exit(statuscode.ProgrammerError)
				}
				logrus.Error(err)
				os.Exit(statuscode.Usage)
			}
			return errors.NotFound{Msg: err.Error()}
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
