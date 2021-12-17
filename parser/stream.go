// Copyright 2020 PingCAP, Inc.
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
)

// StreamParser is a parser implementation which parses bytes from
// io.Reader into individual *LogEntry. Users can parse large log files
// on demand without having to read them all into memory at once.
type StreamParser struct {
	scanner *bufio.Scanner
	Line    int
}

// NewStreamParser creates new *StreamParser associated with the io.Reader.
func NewStreamParser(reader io.Reader) *StreamParser {
	return &StreamParser{
		scanner: bufio.NewScanner(reader),
		Line:    0,
	}
}

// Next reads and parses one LogEntry from bufio.Reader on demand.
// This function will return (nil, nil) if the underlying io.Reader returns
// io.EOF in the standard case.
func (p *StreamParser) Next() (*LogEntry, error) {
	for p.scanner.Scan() {
		p.Line++
		tokens := parseLine(p.scanner.Text())
		if len(tokens) == 0 {
			continue
		}
		p := Parser{tokens: tokens}
		log, err := p.Parse()
		if err != nil {
			return nil, err
		}
		return log, nil
	}
	return nil, p.scanner.Err()
}
