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
	"bytes"
	"io"
	"strconv"
	"strings"
	"text/scanner"
	"time"
)

//	LogEntry 	= LogHeader, Message, LogFields
//	LogHeader 	= '[', DateTime, ']', '[', LogLevel, ']', '[', FileLine, ']'
//	DateTime 	= <string>
//	LogLevel	= "DEBUG" | "INFO" | "WARN" | "ERROR" | "FATAL"
//	FileLine	= <string>
//	Message		= '[', <string>, ']'
//	LogFields	= {LogField}
//	LogField	= '[', [<string>], '=', [<string>], ']' | '[', ']'

// TokenType is an enumeration type for the token type.
type TokenType string

const (
	TokenTypeString   TokenType = "<string>"
	TokenTypeLBracket TokenType = "["
	TokenTypeRBracket TokenType = "]"
	TokenTypeEQ       TokenType = "="
)

type Parser struct {
	tokens []string
}

// Input: "[", "2021/12/13 20:41:00.755 +08:00", "]", "[", "INFO","]", "[", "store.go:68", "]", "[", "new store", "]", "[", "path", "=", "unistore:///tmp/tidb", "]"
// Output:
//    LogEntry: {
//	    Header: {
//		  Datetime: "2021/12/13 20:41:00.755 +08:00",
//        Level: "INFO",
//		  File: "store.go",
//		  Line: 68,
//	    },
//      Message: "new store",
//      Fields: [{
//	      Name: "path",
//        Value: "unistore:///tmp/tidb"
//      }]
//    }
func (p *Parser) Parse() (*LogEntry, error) {
	datetime, err := p.parseDateTime()
	if err != nil {
		return nil, err
	}
	level, err := p.parseLogLevel()
	if err != nil {
		return nil, err
	}
	file, line, err := p.parseFileLine()
	if err != nil {
		return nil, err
	}
	message, err := p.parseMessage()
	if err != nil {
		return nil, err
	}
	fields := []LogField{}
	for len(p.tokens) > 0 {
		field, err := p.parseLogField()
		if err != nil {
			return nil, err
		}
		if field == nil {
			continue
		}
		fields = append(fields, *field)
	}
	return &LogEntry{
		Header: LogHeader{
			DateTime: datetime,
			Level:    level,
			File:     file,
			Line:     line,
		},
		Message: message,
		Fields:  fields,
	}, nil
}

func (p *Parser) parseDateTime() (time.Time, error) {
	if _, err := p.expect(TokenTypeLBracket); err != nil {
		return time.Time{}, err
	}
	tok, err := p.expect(TokenTypeString)
	if err != nil {
		return time.Time{}, err
	}
	datetime, err := time.Parse(TiDBTimeFormat, tok)
	if err != nil {
		return time.Time{}, err
	}
	_, err = p.expect(TokenTypeRBracket)
	if err != nil {
		return time.Time{}, err
	}
	return datetime, nil
}

func (p *Parser) parseLogLevel() (LogLevel, error) {
	_, err := p.expect(TokenTypeLBracket)
	if err != nil {
		return "", err
	}
	tok, err := p.expect(TokenTypeString)
	if err != nil {
		return "", err
	}
	_, err = p.expect(TokenTypeRBracket)
	if err != nil {
		return "", err
	}
	switch tok {
	case string(LogLevelDebug):
		return LogLevelDebug, nil
	case string(LogLevelInfo):
		return LogLevelInfo, nil
	case string(LogLevelWarn):
		return LogLevelWarn, nil
	case string(LogLevelError):
		return LogLevelError, nil
	case string(LogLevelFatal):
		return LogLevelFatal, nil
	default:
		return "", &UnexpectedTokenError{
			ExpectedToken: "LogLevel",
			GotToken:      tok,
		}
	}
}

func (p *Parser) parseFileLine() (string, uint, error) {
	_, err := p.expect(TokenTypeLBracket)
	if err != nil {
		return "", 0, err
	}
	tok, err := p.expect(TokenTypeString)
	if err != nil {
		return "", 0, err
	}
	_, err = p.expect(TokenTypeRBracket)
	if err != nil {
		return "", 0, err
	}
	if tok == "<unknown>" {
		return tok, 0, nil
	}
	xs := strings.Split(tok, ":")
	if len(xs) != 2 {
		return "", 0, &UnexpectedTokenError{
			ExpectedToken: ":",
			GotToken:      "]",
		}
	}
	line, err := strconv.Atoi(xs[1])
	if err != nil {
		return "", 0, err
	}
	return xs[0], uint(line), nil
}

func (p *Parser) parseMessage() (string, error) {
	_, err := p.expect(TokenTypeLBracket)
	if err != nil {
		return "", err
	}
	tok, err := p.expect(TokenTypeString)
	if err != nil {
		return "", err
	}
	_, err = p.expect(TokenTypeRBracket)
	if err != nil {
		return "", err
	}
	return tok, nil
}

func (p *Parser) parseLogField() (*LogField, error) {
	name := ""
	value := ""

	_, err := p.expect(TokenTypeLBracket)
	if err != nil {
		return nil, err
	}

	// peek name
	tok, err := p.peek()
	if err != nil {
		return nil, err
	}
	p.skip()
	if tok == string(TokenTypeRBracket) {
		// "[]", empty log filed, returns empty field
		return nil, nil
	} else if tok != string(TokenTypeEQ) {
		// "[name=xxx]""
		if name, err = unquote(tok); err != nil {
			return nil, err
		}
		if _, err = p.expect(TokenTypeEQ); err != nil {
			return nil, err
		}
	}

	// peek value
	tok, err = p.peek()
	if err != nil {
		return nil, err
	}
	p.skip()
	if tok != string(TokenTypeRBracket) {
		if value, err = unquote(tok); err != nil {
			return nil, err
		}
		_, err = p.expect(TokenTypeRBracket)
		if err != nil {
			return nil, err
		}
	}

	return &LogField{
		Name:  name,
		Value: value,
	}, nil
}

// expect check if p.tokens[0] is expected token and pop it
func (p *Parser) expect(token TokenType) (string, error) {
	if len(p.tokens) == 0 {
		return "", &UnexpectedEOLError{
			ExpectedToken: string(token),
		}
	}
	top := p.tokens[0]
	p.tokens = p.tokens[1:]
	var err error
	if token == TokenTypeString {
		if top, err = unquote(top); err != nil {
			return "", err
		}
	} else if top != string(token) {
		return "", &UnexpectedTokenError{
			ExpectedToken: string(token),
			GotToken:      top,
		}
	}
	return top, nil
}

func (p *Parser) peek() (string, error) {
	if len(p.tokens) == 0 {
		return "", &UnexpectedEOLError{
			ExpectedToken: "<token>",
		}
	}
	return p.tokens[0], nil
}

func (p *Parser) skip() {
	p.tokens = p.tokens[1:]
}

// ParseFromBytes parses a byte slice as *LogEntry slice.
func ParseFromBytes(r []byte) ([]*LogEntry, error) {
	return ParseFromReader(bytes.NewReader(r))
}

// ParseFromString parses a string as *LogEntry slice.
func ParseFromString(r string) ([]*LogEntry, error) {
	return ParseFromReader(strings.NewReader(r))
}

// ParseFromReader parses a byte stream from io.Reader as *LogEntry slice.
// The function continues to run until the reader returns io.EOF.
func ParseFromReader(lr io.Reader) ([]*LogEntry, error) {
	logs := []*LogEntry{}
	sc := bufio.NewScanner(lr)
	for sc.Scan() {
		tokens := parseLine(sc.Text())
		if len(tokens) == 0 {
			continue
		}
		p := Parser{tokens: tokens}
		log, err := p.Parse()
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}
	return logs, nil
}

// Input: [2021/12/13 20:41:00.755 +08:00] [INFO] [store.go:68] ["new store"] [path=unistore:///tmp/tidb]
// Output: "[", "2021/12/13 20:41:00.755 +08:00", "]", "[", "INFO","]", "[", "new store", "]", "[", "path", "=", "unistore:///tmp/tidb", "]"
func parseLine(line string) []string {
	s := scanner.Scanner{}
	s.Init(strings.NewReader(line))

	s.Error = func(s *scanner.Scanner, msg string) {}
	s.IsIdentRune = func(ch rune, i int) bool {
		r := !strings.ContainsRune(`[=]"`, ch) && ch != scanner.EOF
		return r
	}

	xs := []string{}
	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		xs = append(xs, s.TokenText())
	}
	return xs
}

func unquote(s string) (string, error) {
	if len(s) == 0 || s[0] != '"' {
		return s, nil
	}
	return strconv.Unquote(s)
}
