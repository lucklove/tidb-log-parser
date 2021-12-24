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

package event

import (
	"testing"

	"github.com/pingcap-inc/tidb-log-parser/parser"
	"github.com/stretchr/testify/assert"
)

func TestGetLogEventID(t *testing.T) {
	em, err := NewEventManager()
	assert.Nil(t, err)
	id := em.GetLogEventID(&parser.LogEntry{
		Header:  parser.LogHeader{Level: parser.LogLevelInfo},
		Message: "send schedule command",
		Fields: []parser.LogField{{
			Name:  "region-id",
			Value: "xxx",
		}, {
			Name:  "step",
			Value: "xxx",
		}, {
			Name:  "source",
			Value: "xxx",
		}},
	})
	assert.Equal(t, uint(30001), id)

	id = em.GetLogEventID(&parser.LogEntry{
		Header:  parser.LogHeader{Level: parser.LogLevelInfo},
		Message: "TiFlash found 1 stale regions. Only first 1 regions will be logged if the log level is higher than Debug",
		Fields:  []parser.LogField{},
	})
	assert.Equal(t, uint(10112), id)

	id = em.GetLogEventID(&parser.LogEntry{
		Header:  parser.LogHeader{Level: parser.LogLevelInfo},
		Message: "send schedule command",
		Fields:  []parser.LogField{},
	})
	assert.Equal(t, uint(0), id)
}

func TestGussLogEventID(t *testing.T) {
	em, err := NewEventManager()
	assert.Nil(t, err)

	ids := em.GuessLogEventID(&parser.LogEntry{
		Header:  parser.LogHeader{Level: parser.LogLevelInfo},
		Message: "send schedule command",
	}, 1)
	assert.Equal(t, 1, len(ids))
	assert.Equal(t, uint(30001), ids[0])

	ids = em.GuessLogEventID(&parser.LogEntry{
		Header:  parser.LogHeader{Level: parser.LogLevelInfo},
		Message: "send command",
	}, 1)
	assert.Equal(t, 1, len(ids))
	assert.Equal(t, uint(30001), ids[0])

	ids = em.GuessLogEventID(&parser.LogEntry{
		Message: "ddl worker",
	}, 1)
	assert.Equal(t, 1, len(ids))
	assert.Equal(t, uint(10001), ids[0])
}
