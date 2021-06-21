package command

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v2"
)

type CommandSetArgument struct {
	From struct {
		SubCommand string `yaml:"subcommand" json:"subcommand"`
	} `yaml:"from" json:"from"`
	Values []string
}

func (csa *CommandSetArgument) Resolve() ([]string, error) {
	values := []string{}
	if csa.From.SubCommand != "" {
		cmd := exec.Command("/Users/rob/Dev/nidito/milpa/src/milpa", strings.Split(csa.From.SubCommand, " ")...)
		out, err := cmd.Output()
		if err != nil {
			logrus.Error(cmd.CombinedOutput())
			return values, err
		}

		values = strings.Split(strings.TrimSuffix(string(out), "\n"), "\n")
	} else if len(csa.Values) > 0 {
		return csa.Values, nil
	}

	return values, nil
}

type CommandArgument struct {
	Name        string              `yaml:"name" json:"name"`
	Description string              `yaml:"description" json:"description"`
	Set         *CommandSetArgument `yaml:"set" json:"set,omitempty"`
	Variadic    bool                `yaml:"variadic" json:"variadic"`
	Required    bool                `yaml:"required" json:"required"`
}

func (cmdarg *CommandArgument) Validates() bool {
	return cmdarg.Set != nil
}

type CommandOption struct {
	ShortName   string      `yaml:"short-name" json:"short-name"`
	Type        string      `yaml:"type" json:"type"`
	Description string      `yaml:"description" json:"description"`
	Default     interface{} `yaml:"default" json:"default"`
}

type Command struct {
	Meta         CommandMeta              `json:"_meta"`
	Summary      string                   `yaml:"summary" json:"summary"`
	Description  string                   `yaml:"description" json:"description"`
	Arguments    []CommandArgument        `yaml:"arguments" json:"arguments"`
	Options      map[string]CommandOption `yaml:"options" json:"options"`
	runtimeFlags *pflag.FlagSet
}

type CommandMeta struct {
	Path    string   `json:"path"`
	Package string   `json:"package"`
	Name    []string `json:"name"`
	Kind    string   `json:"kind"`
}

func New(path string, spec string, pkg string, kind string) (*Command, error) {
	contents, err := ioutil.ReadFile(spec)
	if err != nil {
		return nil, err
	}

	commandPath := strings.SplitN(path, fmt.Sprintf("%s/.milpa/commands/", pkg), 2)[1]
	commandName := strings.Split(strings.TrimSuffix(commandPath, ".sh"), "/")

	cmd, err := parseCommand(contents, CommandMeta{
		Path:    path,
		Package: pkg,
		Name:    commandName,
		Kind:    kind,
	})
	return cmd, err
}

func parseCommand(yamlBytes []byte, meta CommandMeta) (*Command, error) {
	cmd := &Command{
		Meta:      meta,
		Arguments: []CommandArgument{},
		Options:   map[string]CommandOption{},
	}
	err := yaml.Unmarshal(yamlBytes, cmd)

	if err != nil {
		return nil, err
	}

	return cmd, nil
}

func (cmd *Command) FullName() string {
	return strings.Join(cmd.Meta.Name, " ")
}

func (cmd *Command) ParseArgs(args []string) ([]string, error) {
	logrus.Debugf("Parsing args: %s", args)
	fs := pflag.NewFlagSet(strings.Join(cmd.Meta.Name, " "), pflag.ContinueOnError)
	fs.SortFlags = false
	fs.Usage = func() {}

	for name, opt := range cmd.Options {
		switch opt.Type {
		case "boolean":
			def := false
			if opt.Default != nil {
				def = opt.Default.(bool)
			}
			fs.Bool(name, def, opt.Description)
		case "string":
			def := ""
			if opt.Default != nil {
				def = opt.Default.(string)
			}
			fs.String(name, def, opt.Description)
		default:
			return nil, fmt.Errorf("unknown option type: %s", opt.Type)
		}
	}

	cmd.runtimeFlags = fs
	err := fs.Parse(args)

	if err != nil {
		return nil, err
	}

	for name, opt := range cmd.Options {
		for idx, v := range args {
			modifier := 1
			if opt.Type == "string" {
				modifier = 2
			}

			if v == "--"+name {
				logrus.Debugf("removing %d args, for known flag: %s", modifier, name)
				args = append(args[:idx], args[idx+modifier:]...)
			}

			if opt.ShortName != "" && v == "-"+opt.ShortName {
				args = append(args[:idx], args[idx+modifier:]...)
			}
		}
	}

	return args, nil
}
