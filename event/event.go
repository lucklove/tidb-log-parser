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
	"regexp"
	"sort"
	"strings"

	"github.com/lucklove/tidb-log-parser/parser"
	"github.com/lucklove/tidb-log-parser/utils"
)

// EventManager is responsible for allocate event id for every LogEntry
type EventManager struct {
	// the key is message, the value is an array in case the key is not unique
	msgRule map[string][]*Rule

	// the key is id of the rule
	idRule map[uint][]*Rule

	// regex map
	msgRegex map[string]*regexp.Regexp
}

func NewEventManager(tps ...ComponentType) (*EventManager, error) {
	rs, err := loadRule(tps...)
	if err != nil {
		return nil, err
	}
	msgRule := make(map[string][]*Rule)
	idRule := make(map[uint][]*Rule)
	msgRegex := make(map[string]*regexp.Regexp)
	for _, r := range rs {
		if r.MessageMode() == MessageModeRegex {
			regex, err := regexp.Compile(r.Patterns.Message)
			if err != nil {
				return nil, err
			}
			msgRegex[r.Patterns.Message] = regex
		}
		msgRule[r.Patterns.Message] = append(msgRule[r.Patterns.Message], r)
		idRule[r.ID] = append(idRule[r.ID], r)
	}
	return &EventManager{msgRule, idRule, msgRegex}, nil
}

// GetRuleByID return rules with specified id
func (em *EventManager) GetRulesByEventID(id uint) []*Rule {
	return em.idRule[id]
}

// GetRuleByLog returns the rule matched the log
func (em *EventManager) GetRuleByLog(l *parser.LogEntry) *Rule {
	if rs, ok := em.msgRule[l.Message]; ok {
		if r := em.findMatchedRule(l, rs); r != nil {
			return r
		}
	}

	rules := []*Rule{}
	for msg := range em.msgRule {
		if !strings.Contains(l.Message, msg) {
			continue
		}
		for _, r := range em.msgRule[msg] {
			if r.MessageMode() != MessageModeSubstr {
				continue
			}
			rules = append(rules, r)
		}
	}

	for msg := range em.msgRegex {
		rules = append(rules, em.msgRule[msg]...)
	}

	return em.findMatchedRule(l, rules)
}

// GetLogEventID scan the event conversion rules
// to find the ID for a LogEntry
func (em *EventManager) GetLogEventID(l *parser.LogEntry) uint {
	r := em.GetRuleByLog(l)
	if r == nil {
		return 0
	}
	return r.ID
}

// GuessLogEventID scan the event conversion rules
// and try to find the topN most likely event ids
func (em *EventManager) GuessLogEventID(l *parser.LogEntry, n int) []uint {
	ids := []uint{}
	if n == 0 {
		return ids
	}

	msgs := []string{}
	for msg := range em.msgRule {
		msgs = append(msgs, msg)
	}
	sort.Sort(&stringSorter{l.Message, msgs})

LOOP_MSG:
	for _, msg := range msgs {
		rs := map[string]*Rule{}
		ss := []string{}
		for _, rule := range em.msgRule[msg] {
			k := getKeyFromRule(rule)
			rs[k] = rule
			ss = append(ss, k)
			sort.Sort(&stringSorter{getKeyFromLog(l), ss})
		}
		for _, s := range ss {
			ids = append(ids, rs[s].ID)
			if len(ids) == n {
				break LOOP_MSG
			}
		}
	}

	return ids
}

func (em *EventManager) findMatchedRule(l *parser.LogEntry, rules []*Rule) *Rule {
RULE_LOOP:
	for _, r := range rules {
		if r.MessageMode() == MessageModeRegex {
			if !em.msgRegex[r.Patterns.Message].MatchString(l.Message) {
				continue RULE_LOOP
			}
		} else if r.MessageMode() == MessageModeSubstr {
			if !strings.Contains(l.Message, r.Patterns.Message) {
				continue RULE_LOOP
			}
		} else if l.Message != r.Patterns.Message {
			continue RULE_LOOP
		}
		if r.Patterns.Level != string(l.Header.Level) {
			continue RULE_LOOP
		}

		fns := []string{}
		for _, f := range l.Fields {
			fns = append(fns, f.Name)
		}
		if len(utils.NewStringSet(r.Patterns.Fields...).Difference(utils.NewStringSet(fns...))) > 0 {
			continue RULE_LOOP
		}
		return r
	}
	return nil
}

func getKeyFromLog(l *parser.LogEntry) string {
	xs := []string{string(l.Header.Level)}
	for _, f := range l.Fields {
		xs = append(xs, f.Name)
	}
	return strings.Join(xs, "#")
}

func getKeyFromRule(rule *Rule) string {
	xs := []string{rule.Patterns.Level}
	xs = append(xs, rule.Patterns.Fields...)
	return strings.Join(xs, "#")
}
