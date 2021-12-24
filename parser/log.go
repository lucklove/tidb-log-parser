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
	"time"
)

// LogLevel is an enumeration type for the log level.
type LogLevel string

const (
	LogLevelDebug LogLevel = "DEBUG"
	LogLevelInfo  LogLevel = "INFO"
	LogLevelWarn  LogLevel = "WARN"
	LogLevelError LogLevel = "ERROR"
	LogLevelFatal LogLevel = "FATAL"

	TiDBTimeFormat string = "2006/01/02 15:04:05.000 -07:00"
)

// LogHeader defines the header of one log.
type LogHeader struct {
	DateTime time.Time
	Level    LogLevel
	File     string
	Line     uint
}

// LogField defines one k/v field of one log.
type LogField struct {
	Name  string
	Value string
}

// LogEntry defines an entire log entry.
type LogEntry struct {
	Header  LogHeader
	Message string
	Fields  []LogField // TODO: considering hashmap
}
