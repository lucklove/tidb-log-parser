package main

import (
	"fmt"
	"os"

	"github.com/lucklove/tidb-log-parser/event"
	"github.com/lucklove/tidb-log-parser/parser"
	"github.com/spf13/cobra"
)

func newExportCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "export",
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
				eid := em.GetLogEventID(log)
				if eid == 0 {
					continue
				}
				fmt.Printf("%d,%d\n", log.Header.DateTime.Unix(), eid)
			}
			return nil
		},
	}

	return cmd
}
