/*-
 * Copyright 2015 Grammarly, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package build2

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"rocker/parser"
	"rocker/template"
	"strings"
)

type Rockerfile struct {
	Name    string
	Source  string
	Content string
	Vars    template.Vars
	Funs    template.Funs

	rootNode *parser.Node
}

func NewRockerfileFromFile(name string, vars template.Vars, funs template.Funs) (r *Rockerfile, err error) {
	fd, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	return NewRockerfile(name, fd, vars, funs)
}

func NewRockerfile(name string, in io.Reader, vars template.Vars, funs template.Funs) (r *Rockerfile, err error) {
	r = &Rockerfile{
		Name: name,
		Vars: vars,
		Funs: funs,
	}

	var (
		source  []byte
		content *bytes.Buffer
	)

	if source, err = ioutil.ReadAll(in); err != nil {
		return nil, fmt.Errorf("Failed to read Rockerfile %s, error: %s", name, err)
	}

	r.Source = string(source)

	if content, err = template.Process(name, bytes.NewReader(source), vars, funs); err != nil {
		return nil, err
	}

	r.Content = content.String()

	if r.rootNode, err = parser.Parse(content); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *Rockerfile) Commands() []ConfigCommand {
	commands := []ConfigCommand{}

	for i := 0; i < len(r.rootNode.Children); i++ {
		node := r.rootNode.Children[i]

		cfg := ConfigCommand{
			name:     node.Value,
			attrs:    node.Attributes,
			original: node.Original,
			args:     []string{},
			flags:    parseFlags(node.Flags),
		}

		// fill in args and substitute vars
		for n := node.Next; n != nil; n = n.Next {
			cfg.args = append(cfg.args, n.Value)
		}

		commands = append(commands, cfg)
	}

	return commands
}

func parseFlags(flags []string) map[string]string {
	result := make(map[string]string)
	for _, flag := range flags {
		key := flag[2:]
		value := ""

		index := strings.Index(key, "=")
		if index >= 0 {
			value = key[index+1:]
			key = key[:index]
		}

		result[key] = value
	}
	return result
}