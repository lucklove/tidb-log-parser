// Copyright 2021 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/lucklove/tidb-log-parser/event"
	"github.com/lucklove/tidb-log-parser/parser"
	"github.com/spf13/cobra"
)

func newCheckCommand() *cobra.Command {
	withoutTime := false
	cmd := &cobra.Command{
		Use: "check <component>",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.Help()
			}
			comp, err := event.GetComponentType(args[0])
			if err != nil {
				return err
			}
			p := parser.NewStreamParser(os.Stdin)
			if withoutTime {
				p = p.WithoutTime()
			}
			em, err := event.NewEventManager(comp)
			assert(err)

			for {
				log, err := p.Next()
				if log == nil && err == nil {
					break
				}
				if log == nil || err != nil {
					fmt.Println(err)
					continue
				}
				if ignore(log) {
					continue
				}
				if em.GetLogEventID(log) == 0 {
					xs := em.GuessLogEventID(log, 1)
					rs := em.GetRulesByEventID(xs[0])
					fs := []string{}
					for _, f := range log.Fields {
						fs = append(fs, f.Name)
					}
					r := event.Rule{
						ID:   xs[0],
						Name: rs[0].Name,
						Patterns: event.RulePattern{
							Level:   string(log.Header.Level),
							Message: log.Message,
							Fields:  fs,
						},
					}
					err = toml.NewEncoder(os.Stdout).Encode(struct {
						Rule []event.Rule `toml:"rule"`
					}{Rule: []event.Rule{r}})
					assert(err)
				}
			}
			return nil
		},
	}

	cmd.Flags().BoolVarP(&withoutTime, "without-time", "", false, "if every line doesn't contains the time header")
	return cmd
}

func assert(err error) {
	if err != nil {
		panic(err)
	}
}

func ignore(log *parser.LogEntry) bool {
	ignoreFiles := []string{}
	ignoreMessages := []string{}

	for _, f := range ignoreFiles {
		if f == log.Header.File {
			return true
		}
	}
	for _, m := range ignoreMessages {
		if strings.Contains(log.Message, m) {
			return true
		}
	}
	return false
}
