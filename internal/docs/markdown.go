// SPDX-License-Identifier: Apache-2.0
// Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>
package docs

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"sort"
	"strings"

	"git.rob.mx/nidito/chinampa/pkg/command"
	"git.rob.mx/nidito/chinampa/pkg/tree"
	"github.com/sirupsen/logrus"
	"github.com/unrob/milpa/internal/lookup"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
	"gopkg.in/yaml.v3"
)

type FrontMatter struct {
	Weight      int
	Description string
	Title       string
}

func frontMatter(contents []byte) (*FrontMatter, []byte) {
	frontmatterSep := []byte("---\n")
	if len(contents) > 4 && strings.Contains(string(contents[0:5]), string(frontmatterSep)) {
		parts := bytes.SplitN(contents, frontmatterSep, 3)
		res := &FrontMatter{}
		if err := yaml.Unmarshal(parts[1], &res); err == nil {
			return res, parts[2]
		}
	}

	return nil, contents
}

type Page struct {
	Name     string
	Path     string
	Weight   int
	Children *Pages
}
type Pages []*Page

var _ sort.Interface = Pages{}

func (p Pages) Len() int { return len(p) }
func (p Pages) Less(i, j int) bool {
	if p[i].Weight != p[j].Weight {
		return p[i].Weight < p[j].Weight
	}
	return p[i].Path < p[j].Path
}
func (p Pages) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

func (p Page) Sort() {
	for _, p := range *p.Children {
		p.Sort()
	}
	sort.Sort(p.Children)
}

type Entry struct {
	ID      string
	Title   template.HTML
	Entries Entries
}

func (i *Entry) Append() *Entry {
	child := &Entry{}
	i.Entries = append(i.Entries, child)
	return child
}

func (i *Entry) Last() *Entry {
	if len(i.Entries) > 0 {
		return i.Entries[len(i.Entries)-1]
	}
	return i.Append()
}

type Entries []*Entry

type tocTransformer struct {
	Entries *Entries
}

var _ parser.ASTTransformer = &tocTransformer{}

func (t *tocTransformer) Transform(doc *ast.Document, reader text.Reader, pctx parser.Context) {
	root := &Entry{Entries: Entries{}}

	current := []*Entry{root}
	src := reader.Source()
	err := ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		heading, ok := n.(*ast.Heading)
		idI, hasID := n.AttributeString("id")
		if !ok || !hasID {
			return ast.WalkContinue, nil
		}

		for len(current) < heading.Level {
			parent := current[len(current)-1]
			current = append(current, parent.Last())
		}

		for len(current) > heading.Level {
			current = current[:len(current)-1]
		}

		parent := current[len(current)-1]
		target := parent.Last()
		if len(target.Title) > 0 || len(target.Entries) > 0 {
			target = parent.Append()
		}

		rdr := goldmark.DefaultRenderer()
		data := bytes.Buffer{}

		err := ast.Walk(heading, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
			_, isHeading := n.(*ast.Heading)
			if !entering || isHeading {
				return ast.WalkContinue, nil
			}

			if err := rdr.Render(&data, src, n); err != nil {
				logrus.Errorf("could not render: %s", err)
				return ast.WalkStop, nil
			}
			return ast.WalkSkipChildren, nil
		})
		if err != nil {
			log.Errorf("error walking ast: %s", err)
		}

		target.Title = template.HTML(data.String()) // nolint:gosec

		id, _ := idI.([]byte)
		target.ID = string(id)

		return ast.WalkSkipChildren, nil
	})

	if err != nil {
		log.Errorf("error walking ast: %s", err)
	}

	if len(root.Entries) > 0 {
		t.Entries = &root.Entries[0].Entries
	} else {
		t.Entries = &Entries{}
	}
}

type milpaRenderer struct {
	html.Config
	me *milpaExtension
}

func (r *milpaRenderer) SetOption(name renderer.OptionName, value any) {
	r.Config.SetOption(name, value)
}

func (r *milpaRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindHeading, r.renderHeading)
}

func (r *milpaRenderer) renderHeading(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	hn := node.(*ast.Heading)
	text := string(hn.Text(source))

	slug, ok := node.AttributeString("id")
	if !ok {
		log.Errorf("could not get id for %s", text)
		return ast.WalkSkipChildren, nil
	}

	if entering {
		node.SetAttribute([]byte("id"), slug)
		log.Tracef("entering header %s (%d)", node.Text(source), node.ChildCount())
		_, err := w.WriteString(fmt.Sprintf(`<div class="content-header-wrapper"><h%d class="content-header" id="%s">`, hn.Level, slug))
		if err != nil {
			log.Errorf("Error writing header: %s", err)
		}
		return ast.WalkContinue, nil
	}

	log.Tracef("closing header %s (%d)", text, node.ChildCount())

	_, err := w.WriteString(fmt.Sprintf(`</h%d> <a aria-hidden="true" class="heading-anchor" href="#%s" tabindex="-1">#</a></div>`, hn.Level, slug))
	if err != nil {
		log.Errorf("Error writing header: %s", err)
	}

	return ast.WalkContinue, nil
}

type milpaExtension struct {
	TOC *tocTransformer
}

// Extend implements goldmark.Extender.
func (me *milpaExtension) Extend(m goldmark.Markdown) {
	me.TOC = &tocTransformer{Entries: &Entries{}}
	m.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(&milpaRenderer{
				me: me, Config: html.Config{
					Writer: html.DefaultWriter,
				}}, 999),
		),
		html.WithUnsafe(),
	)
	m.Parser().AddOptions(
		parser.WithAutoHeadingID(),
		parser.WithASTTransformers(
			util.Prioritized(me.TOC, 100),
		),
	)
}

func commandSerializer(pageTree *Page, names *[]string) func(cmd *command.Command) error {
	return func(cmd *command.Command) error {
		current := pageTree
		path := cmd.Path[1:]
		log.Tracef("creating tree for %s, current: %d", path, len(*current.Children))
		for idx, c := range path {
			if len(path)-1 == idx {
				*current.Children = append(*current.Children, &Page{
					Name:     cmd.Name(),
					Path:     strings.Join(path, "/"),
					Children: &Pages{},
				})
				*names = append(*names, strings.Join(cmd.Path[1:], " "))
				log.Tracef("inserted %s at %s (%d)", cmd.Name(), current.Path, len(*current.Children))
				return nil
			}

			var newCurrent Page
			found := false
			for _, child := range *current.Children {
				if child.Name == c {
					found = true
					newCurrent = *child
					log.Tracef("found %s in %s", child.Name, current.Path)
					break
				}
			}

			if !found {
				newCurrent = Page{
					Name:     c,
					Path:     strings.Join(path[:idx], "/"),
					Children: &Pages{},
				}
				*current.Children = append(*current.Children, &newCurrent)
				log.Tracef("adding sub-level %s to %s", newCurrent.Name, current.Path)
			}

			current = &newCurrent
		}
		return nil
	}
}

func docsSerializer(pageTree *Page, names *[]string) error {
	allDocs, err := lookup.AllDocs()
	if err != nil {
		return err
	}
	sort.Strings(allDocs)
	log.Infof("Found %d docs", len(allDocs))

	for _, doc := range allDocs {
		current := pageTree
		d := strings.TrimSuffix(strings.SplitN(doc, ".milpa/docs/", 2)[1], ".md")
		path := append([]string{"help", "docs"}, strings.Split(d, "/")...)
		if path[len(path)-1] == "index" {
			path = path[0 : len(path)-1]
		}

		log.Tracef("creating tree for %s, path: %s", path, doc)

		for idx, c := range path {
			if len(path)-1 == idx {
				weight := 999
				contents, err := os.ReadFile(doc)
				if err == nil {
					if fm, _ := frontMatter(contents); fm != nil {
						weight = fm.Weight
					}
				}

				found := false
				for _, child := range *current.Children {
					if child.Name == c {
						found = true
						child.Name = c
						child.Path = strings.Join(path, "/")
						child.Weight = weight
						log.Tracef("reset child %s of %s / %s", c, current.Path, path)
						break
					}
				}

				if !found {
					*current.Children = append(*current.Children, &Page{
						Name:     c,
						Path:     strings.Join(path, "/"),
						Weight:   weight,
						Children: &Pages{},
					})
					log.Tracef("added %s to %s / %s", c, current.Path, path)
				}
				*names = append(*names, strings.Join(path, " "))

				break
			}

			var found *Page
			if current.Children == nil {
				current.Children = &Pages{}
			}

			for _, child := range *current.Children {
				if child.Name == c {
					found = child
					log.Tracef("found %s in %s", child.Name, current.Path)
					break
				}
			}

			if found == nil {
				found = &Page{
					Name:     c,
					Path:     strings.Join(path[:idx+1], "/"),
					Children: &Pages{},
				}
				*current.Children = append(*current.Children, found)
				log.Tracef("adding sub-level %s to %s", found.Name, current.Path)
			}
			current = found
		}
	}

	return nil
}

func buildSiteTree() (*Page, string, error) {
	names := []string{"", "help", "help docs"}
	pageTree := &Page{
		Path: "milpa",
		Name: "milpa",
		Children: &Pages{
			&Page{
				Name:     "help",
				Path:     "help",
				Children: &Pages{},
			},
		},
	}

	// Add commands to tree
	tree.Build(command.Root.Cobra, 20)
	_, err := tree.Serialize(func(t interface{}) ([]byte, error) {
		tree := t.(*tree.CommandTree)
		err := tree.Traverse(commandSerializer(pageTree, &names))
		log.Infof("Found %d commands", len(names))
		return nil, err
	})
	if err != nil {
		log.Errorf("could not build command tree: %s", err)
		return nil, "", err
	}

	// Add docs to tree
	if err := docsSerializer(pageTree, &names); err != nil {
		log.Errorf("could not docs tree: %s", err)
		return nil, "", err
	}

	pageTree.Sort()
	sort.Strings(names)
	return pageTree, strings.Join(names, "|"), nil
}
