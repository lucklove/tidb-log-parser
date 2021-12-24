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
	"math"
	"os"
	"path"
	"sort"

	"github.com/pingcap-inc/tidb-log-parser/event"
	"github.com/pingcap-inc/tidb-log-parser/parser"
	"github.com/pingcap-inc/tidb-log-parser/store"
	"github.com/spf13/cobra"
)

type BatchDiager struct {
	m map[uint]uint
	c uint
}

func (d *BatchDiager) Consume(id uint) {
	d.m[id]++
	d.c++
}

func (d *BatchDiager) Produce() []uint {
	xs := []uint{}
	for eid := range d.m {
		xs = append(xs, eid)
	}
	return xs
}

func (d *BatchDiager) Weight(id uint) float64 {
	return float64(d.m[id]) / float64(d.c)
}

func (d *BatchDiager) Count(id uint) uint {
	return d.m[id]
}

func newDiagCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "diag",
		RunE: func(cmd *cobra.Command, args []string) error {
			d := BatchDiager{make(map[uint]uint), 0}

			home, err := os.UserHomeDir()
			assert(err)
			store, err := store.NewSQLiteStorage(path.Join(home, ".tiup/storage/naglfar/log.db"), "tidb")
			assert(err)
			defer store.Close()

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
				if ignore(log) {
					continue
				}
				eid := em.GetLogEventID(log)
				if eid == 0 {
					fmt.Println(log.Message)
					panic("eid should not be zero, please run check command first")
				}
				d.Consume(eid)
			}

			wm := make(map[uint]float64)
			eids := d.Produce()
			for _, eid := range eids {
				ec, err := store.EventCount(eid)
				assert(err)
				lfc, err := store.LogFragmentCount()
				assert(err)
				wm[eid] = d.Weight(eid) * math.Log(float64(lfc)/float64(ec+1))
			}

			sort.Slice(eids, func(i, j int) bool {
				return wm[eids[i]] > wm[eids[j]]
			})

			for _, eid := range eids {
				rs := em.GetRulesByEventID(eid)
				fmt.Printf("%f\t%d\t%s\n", wm[eid], d.Count(eid), rs[0].Name)
			}

			return nil
		},
	}

	return cmd
}
