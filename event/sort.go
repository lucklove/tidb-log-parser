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
	"github.com/hbollon/go-edlib"
)

// sort strings by similarity with key
type stringSorter struct {
	key  string
	strs []string
}

func (s *stringSorter) Len() int {
	return len(s.strs)
}

func (s *stringSorter) Swap(i, j int) {
	s.strs[i], s.strs[j] = s.strs[j], s.strs[i]
}

func (s *stringSorter) Less(i, j int) bool {
	r1, _ := edlib.StringsSimilarity(s.key, s.strs[i], edlib.Levenshtein)
	r2, _ := edlib.StringsSimilarity(s.key, s.strs[j], edlib.Levenshtein)
	return r1 > r2
}
