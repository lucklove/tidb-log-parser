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
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/pingcap-inc/tidb-log-parser/event"
	"github.com/pingcap-inc/tidb-log-parser/parser"
	"github.com/spf13/cobra"
)

func newCheckCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "check",
		RunE: func(cmd *cobra.Command, args []string) error {
			p := parser.NewStreamParser(os.Stdin)
			em, err := event.NewEventManager(event.ComponentTiDB)
			assert(err)

			for {
				log, err := p.Next()
				if log == nil && err == nil {
					break
				}
				if log == nil || err != nil {
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
