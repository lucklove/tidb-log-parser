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
	"path"

	"github.com/lucklove/tidb-log-parser/event"
	"github.com/lucklove/tidb-log-parser/parser"
	"github.com/lucklove/tidb-log-parser/store"
	"github.com/spf13/cobra"
)

func newLearnCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "learn",
		RunE: func(cmd *cobra.Command, args []string) error {
			home, err := os.UserHomeDir()
			assert(err)
			store, err := store.NewSQLiteStorage(path.Join(home, ".tiup/storage/naglfar/log.db"), "tidb")
			assert(err)
			defer store.Close()
			fc, err := store.LogFragmentCount()
			assert(err)
			fid := fc + 1

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
				eid := em.GetLogEventID(log)
				if eid == 0 {
					fmt.Println(log.Message)
					panic("eid should not be zero, please run check command first")
				}
				store.Insert(fid, eid)
			}
			return nil
		},
	}

	return cmd
}
