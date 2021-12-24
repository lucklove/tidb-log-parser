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
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSort(t *testing.T) {
	ss := []string{
		"Copyright 2020 PingCAP, Inc.",
		"Licensed under the Apache License",
		"you may not use this file except in compliance with the License.",
		"Unless required by applicable law or agreed to in writing",
		`distributed under the License is distributed on an "AS IS" BASIS`,
	}
	sort.Sort(&stringSorter{"the license is on some what AS IS basis", ss})
	assert.Equal(t, `distributed under the License is distributed on an "AS IS" BASIS`, ss[0])

	sort.Sort(&stringSorter{"under Apache License", ss})
	assert.Equal(t, "Licensed under the Apache License", ss[0])
}
