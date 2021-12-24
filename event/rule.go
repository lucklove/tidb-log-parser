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
	_ "embed"

	"github.com/BurntSushi/toml"
)

//go:embed tidb.toml
var tidbRuleStr string

//go:embed tikv.toml
var tikvRuleStr string

//go:embed pd.toml
var pdRuleStr string

//go:embed tiflash.toml
var tiflashRuleStr string

type ComponentType int
type MessageModeType int

const (
	ComponentTiDB    ComponentType = iota
	ComponentTiKV    ComponentType = iota
	ComponentPD      ComponentType = iota
	ComponentTiFlash ComponentType = iota

	MessageModeEqual  MessageModeType = iota
	MessageModeRegex  MessageModeType = iota
	MessageModeSubstr MessageModeType = iota
)

// Rule indicates how to convert LogEntry to event
type Rule struct {
	ID       uint        `toml:"id"`
	Name     string      `toml:"name"`
	Patterns RulePattern `toml:"patterns"`
}

// RulePattern is a selector which describle how the LogEntry looks like
type RulePattern struct {
	Level       string   `toml:"level"`
	Message     string   `toml:"message"`
	MessageMode string   `toml:"message_mode"`
	Fields      []string `toml:"fields"`
}

func (r *Rule) MessageMode() MessageModeType {
	switch r.Patterns.MessageMode {
	case "regex":
		return MessageModeRegex
	case "substr":
		return MessageModeSubstr
	default:
		return MessageModeEqual
	}
}

func loadRule(tps ...ComponentType) ([]*Rule, error) {
	if len(tps) == 0 {
		tps = []ComponentType{
			ComponentTiDB, ComponentTiKV, ComponentPD, ComponentTiFlash,
		}
	}

	rules := []*Rule{}
	for _, tp := range tps {
		rs := struct {
			Rule []*Rule `toml:"rule"`
		}{}
		switch tp {
		case ComponentTiDB:
			if _, err := toml.Decode(tidbRuleStr, &rs); err != nil {
				return nil, err
			}
		case ComponentTiKV:
			if _, err := toml.Decode(tikvRuleStr, &rs); err != nil {
				return nil, err
			}
		case ComponentPD:
			if _, err := toml.Decode(pdRuleStr, &rs); err != nil {
				return nil, err
			}
		case ComponentTiFlash:
			if _, err := toml.Decode(tiflashRuleStr, &rs); err != nil {
				return nil, err
			}
		default:
			panic("unreachable")
		}
		rules = append(rules, rs.Rule...)
	}
	return rules, nil
}
