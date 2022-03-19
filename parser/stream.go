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
	"bufio"
	"io"

	"github.com/pingcap/errors"
)

// StreamParser is a parser implementation which parses bytes from
// io.Reader into individual *LogEntry. Users can parse large log files
// on demand without having to read them all into memory at once.
type StreamParser struct {
	scanner     *bufio.Scanner
	Line        int
	withoutTime bool
}

// NewStreamParser creates new *StreamParser associated with the io.Reader.
func NewStreamParser(reader io.Reader) *StreamParser {
	scanner := bufio.NewScanner(reader)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)
	return &StreamParser{
		scanner: scanner,
		Line:    0,
	}
}

func (sp *StreamParser) WithoutTime() *StreamParser {
	sp.withoutTime = true
	return sp
}

// Next reads and parses one LogEntry from bufio.Reader on demand.
// This function will return (nil, nil) if the underlying io.Reader returns
// io.EOF in the standard case.
func (sp *StreamParser) Next() (*LogEntry, error) {
	for sp.scanner.Scan() {
		sp.Line++
		line := sp.scanner.Text()
		if sp.withoutTime {
			line = "[2006/01/02 15:04:05.000 -07:00] " + line
		}
		tokens := parseLine(line)
		if len(tokens) == 0 {
			continue
		}
		p := Parser{tokens: tokens}
		log, err := p.Parse()
		if err != nil {
			return nil, errors.Annotatef(err, "at line %d", sp.Line)
		}
		return log, nil
	}
	sp.Line++
	return nil, errors.Annotatef(sp.scanner.Err(), "at line %d", sp.Line)
}
