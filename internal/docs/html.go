// SPDX-License-Identifier: Apache-2.0
// Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>
package docs

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"git.rob.mx/nidito/chinampa/pkg/command"
	"git.rob.mx/nidito/chinampa/pkg/render"
	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/spf13/cobra"
	_c "github.com/unrob/milpa/internal/constants"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"

	"github.com/yuin/goldmark/extension"
)

//go:embed template.html
var LayoutTemplate []byte

//go:embed static/*
var StaticFiles embed.FS

type TemplateContents struct {
	Base           string
	IsHome         bool
	Permalink      string
	RelPermalink   string
	Content        template.HTML
	Description    string
	Tree           *Page
	TOC            *Entries
	CommandPattern string
}

func FixLinks(contents []byte) []byte {
	fixedLinks := bytes.ReplaceAll(contents, []byte("(/"+_c.RepoDocs), []byte("(/help/docs"))
	fixedLinks = bytes.ReplaceAll(fixedLinks, []byte("(/"+_c.RepoCommands+"/"), []byte("(/"))
	fixedLinks = bytes.ReplaceAll(fixedLinks, []byte("index.md"), []byte(""))
	return bytes.ReplaceAll(fixedLinks, []byte(".md"), []byte("/"))
}

func getHTMLLayout() (*template.Template, error) {
	return template.New("html-help").Funcs(render.TemplateFuncs).Parse(string(LayoutTemplate))
}

func contentsForRequest(comps []string) ([]byte, string, error) {
	var cmd *cobra.Command
	var args []string
	var err error
	if len(comps) == 1 && comps[0] == "help" {
		cmd = command.Root.Cobra
		args = []string{}
	} else {
		cmd, args, err = command.Root.Cobra.Find(comps)
	}

	var helpMD bytes.Buffer
	if err != nil || (cmd == cmd.Root() && len(args) > 0) {
		log.Infof("returning 404: %s, cmd: %s", comps, cmd.Name())
		contents := []byte("# Not found\n\nThat is weird, if you have a second and a github account, [let me know](https://github.com/unRob/milpa/issues/new?labels=docs&title=Page+not+found&template=docs-page-not-found.yml).\n")
		desc := "sub-command not found"
		return contents, desc, fmt.Errorf("not found: %s", comps)
	}

	fuckingDocs := len(args) == 2 && (args[0] == "help" && args[1] == "docs")
	if cmd.Name() == "docs" || fuckingDocs {
		log.Debugf("Rendering docs for args %s", args)
		helpMD.WriteString("\n")
		if len(args) == 0 || fuckingDocs {
			log.Info("Rendering docs main help page")
			cmd.SetOutput(&helpMD)
			if err := cmd.Help(); err != nil {
				return nil, "", fmt.Errorf("error: %s", err)
			}
		} else {
			log.Infof("Rendering docs topic for %s", args)
			data, err := FromQuery(args)
			if err != nil {
				return nil, "", fmt.Errorf("error: %s", err)
			}
			helpMD.Write(data)
		}
	} else {
		log.Infof("Rendering command help for %s, args: %s", cmd.Use, args)
		cmd.SetOutput(&helpMD)

		if err := cmd.Help(); err != nil {
			return nil, "", fmt.Errorf("error: %s", err)
		}
		log.Debugf("Rendered %s bytes for %s", helpMD.String(), cmd.Name())
	}

	desc := cmd.Short
	fm, contents := frontMatter(helpMD.Bytes())
	if fm != nil && fm.Description != "" {
		desc = fm.Description
	}

	return contents, desc, nil
}

func mdToHTML(md []byte) (bytes.Buffer, *Entries, error) {
	var helpHTML bytes.Buffer

	milpaHeadings := milpaExtension{}
	markdown := goldmark.New(
		goldmark.WithExtensions(
			&milpaHeadings,
			extension.GFM,
			highlighting.NewHighlighting(
				highlighting.WithStyle("xcode"),
				highlighting.WithFormatOptions(
					chromahtml.WithClasses(true),
				),
			),
			extension.Table,
			extension.Strikethrough,
		),
	)

	err := markdown.Convert(FixLinks(md), &helpHTML)
	return helpHTML, milpaHeadings.TOC.Entries, err
}

func StaticHandler() http.Handler {
	fs := http.FS(StaticFiles)
	return http.FileServer(fs)
}

func RenderHandler(serverAddr string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, ".ico") {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		prefix := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/"), "/")
		comps := []string{}
		if prefix != "" {
			comps = strings.Split(prefix, "/")
		}

		log.Infof("Handling request for: %s", comps)

		contents, desc, err := contentsForRequest(comps)
		if err != nil {
			if !strings.HasPrefix(err.Error(), "not found:") {
				log.Errorf("could not get contents: %s", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			log.Errorf("404: %s", comps)
			w.WriteHeader(http.StatusNotFound)
		}

		md, toc, err := mdToHTML(contents)
		if err != nil {
			log.Errorf("could convert to html: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		tpl, err := getHTMLLayout()
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		pageTree, cp, err := buildSiteTree()
		if err != nil {
			log.Errorf("could not build site tree %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var pageHTML bytes.Buffer
		err = tpl.Execute(&pageHTML, &TemplateContents{
			Base:           serverAddr,
			IsHome:         len(comps) == 0,
			RelPermalink:   "/" + prefix,
			Permalink:      serverAddr + "/" + prefix,
			Content:        template.HTML(md.String()), // nolint:gosec
			Description:    desc,
			Tree:           pageTree,
			TOC:            toc,
			CommandPattern: cp,
		})

		if err != nil {
			log.Errorf("could not render template: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Add("content-type", "text/html")
		if _, err := w.Write(pageHTML.Bytes()); err != nil {
			log.Errorf("could not write response: %s", err)
		}
	}
}
