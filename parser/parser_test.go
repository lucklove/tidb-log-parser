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

package parser

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseFromString(t *testing.T) {
	tests := map[string][]*LogEntry{
		// normal test
		`[2021/12/13 20:41:00.755 +08:00] [INFO] [main.go:336] ["disable Prometheus push client"]`: {{
			Header: LogHeader{
				DateTime: time.Date(2021, 12, 13, 20, 41, 00, 755000000, time.FixedZone("CST", 3600*8)),
				Level:    "INFO",
				File:     "main.go",
				Line:     336,
			},
			Message: "disable Prometheus push client",
			Fields:  []LogField{},
		}},
		// test empty fields
		`[2021/12/14 14:21:17.639 +08:00] [WARN] [misc.go:446] ["Automatic TLS Certificate creation is disabled"] []`: {{
			Header: LogHeader{
				DateTime: time.Date(2021, 12, 14, 14, 21, 17, 639000000, time.FixedZone("CST", 3600*8)),
				Level:    "WARN",
				File:     "misc.go",
				Line:     446,
			},
			Message: "Automatic TLS Certificate creation is disabled",
			Fields:  []LogField{},
		}},
		// test tiflash log
		`[2021/12/14 11:02:06.826 +08:00] [INFO] [<unknown>] ["IOLimitTuner: limiter 0 write 0 read 0 NOT need to tune."] [thread_id=6]`: {{
			Header: LogHeader{
				DateTime: time.Date(2021, 12, 14, 11, 02, 06, 826000000, time.FixedZone("CST", 3600*8)),
				Level:    "INFO",
				File:     "<unknown>",
				Line:     0,
			},
			Message: "IOLimitTuner: limiter 0 write 0 read 0 NOT need to tune.",
			Fields: []LogField{{
				Name:  "thread_id",
				Value: "6",
			}},
		}},
		// empty field name
		`[2021/12/14 11:02:06.826 +08:00] [INFO] [<unknown>] ["IOLimitTuner: limiter 0 write 0 read 0 NOT need to tune."] [=7]`: {{
			Header: LogHeader{
				DateTime: time.Date(2021, 12, 14, 11, 02, 06, 826000000, time.FixedZone("CST", 3600*8)),
				Level:    "INFO",
				File:     "<unknown>",
				Line:     0,
			},
			Message: "IOLimitTuner: limiter 0 write 0 read 0 NOT need to tune.",
			Fields: []LogField{{
				Name:  "",
				Value: "7",
			}},
		}},
		// empty field value
		`[2021/12/14 11:02:06.826 +08:00] [INFO] [<unknown>] ["IOLimitTuner: limiter 0 write 0 read 0 NOT need to tune."] [thread_id=]`: {{
			Header: LogHeader{
				DateTime: time.Date(2021, 12, 14, 11, 02, 06, 826000000, time.FixedZone("CST", 3600*8)),
				Level:    "INFO",
				File:     "<unknown>",
				Line:     0,
			},
			Message: "IOLimitTuner: limiter 0 write 0 read 0 NOT need to tune.",
			Fields: []LogField{{
				Name:  "thread_id",
				Value: "",
			}},
		}},
		// both name and value are empty
		`[2021/12/14 11:02:06.826 +08:00] [INFO] [<unknown>] ["IOLimitTuner: limiter 0 write 0 read 0 NOT need to tune."] [=]`: {{
			Header: LogHeader{
				DateTime: time.Date(2021, 12, 14, 11, 02, 06, 826000000, time.FixedZone("CST", 3600*8)),
				Level:    "INFO",
				File:     "<unknown>",
				Line:     0,
			},
			Message: "IOLimitTuner: limiter 0 write 0 read 0 NOT need to tune.",
			Fields: []LogField{{
				Name:  "",
				Value: "",
			}},
		}},
		// multiple fileds with same name
		`[2021/12/14 11:02:06.826 +08:00] [INFO] [<unknown>] ["IOLimitTuner: limiter 0 write 0 read 0 NOT need to tune."] [x=1] [x=2]`: {{
			Header: LogHeader{
				DateTime: time.Date(2021, 12, 14, 11, 02, 06, 826000000, time.FixedZone("CST", 3600*8)),
				Level:    "INFO",
				File:     "<unknown>",
				Line:     0,
			},
			Message: "IOLimitTuner: limiter 0 write 0 read 0 NOT need to tune.",
			Fields: []LogField{{
				Name:  "x",
				Value: "1",
			}, {
				Name:  "x",
				Value: "2",
			}},
		}},
		// unquote
		`[2021/12/14 11:02:06.826 +08:00] [INFO] [<unknown>] ["this a a \"quoted\" string, with [bracket] and [=]"]`: {{
			Header: LogHeader{
				DateTime: time.Date(2021, 12, 14, 11, 02, 06, 826000000, time.FixedZone("CST", 3600*8)),
				Level:    "INFO",
				File:     "<unknown>",
				Line:     0,
			},
			Message: `this a a "quoted" string, with [bracket] and [=]`,
			Fields:  []LogField{},
		}},
	}

	for str, log := range tests {
		l, err := ParseFromString(str)
		assert.Nil(t, err)
		assert.Equal(t, l[0].Header.DateTime.Format(time.RFC3339), log[0].Header.DateTime.Format(time.RFC3339))
		n := time.Now()
		l[0].Header.DateTime = n
		log[0].Header.DateTime = n
		assert.Equal(t, l, log)
	}
}
