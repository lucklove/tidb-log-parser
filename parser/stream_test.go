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
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStream(t *testing.T) {
	logtxt := `
	[2021/12/16 17:03:48.696 +08:00] [INFO] [trackerRecorder.go:28] ["Mem Profile Tracker started"]
	[2021/12/16 17:03:48.696 +08:00] [INFO] [printer.go:33] ["Welcome to TiDB."] ["Release Version"=v5.2.0] [Edition=Community] ["Git Commit Hash"=05d2210647d6a1503a8d772477e43b14a024f609] ["Git Branch"=heads/refs/tags/v5.2.0] ["UTC Build Time"="2021-08-27 05:54:47"] [GoVersion=go1.16.5] ["Race Enabled"=false] ["Check Table Before Drop"=false] ["TiKV Min Version"=v3.0.0-60965b006877ca7234adaced7890d7b029ed1306]
	<invalid line>
	[2021/12/16 17:03:48.698 +08:00] [INFO] [printer.go:47] ["loaded config"] [config="{\"host\":\"0.0.0.0\",\"advertise-address\":\"10.71.200.107\",\"port\":4000,\"cors\":\"\",\"store\":\"unistore\",\"path\":\"/tmp/tidb\",\"socket\":\"\",\"lease\":\"45s\",\"run-ddl\":true,\"split-table\":true,\"token-limit\":1000,\"oom-use-tmp-storage\":true,\"tmp-storage-path\":\"/var/folders/pb/ssv6rnq92pn_l3c9gr9rd1g80000gn/T/501_tidb/MC4wLjAuMDo0MDAwLzAuMC4wLjA6MTAwODA=/tmp-storage\",\"oom-action\":\"cancel\",\"mem-quota-query\":1073741824,\"tmp-storage-quota\":-1,\"enable-batch-dml\":false,\"lower-case-table-names\":2,\"server-version\":\"\",\"log\":{\"level\":\"info\",\"format\":\"text\",\"disable-timestamp\":null,\"enable-timestamp\":null,\"disable-error-stack\":null,\"enable-error-stack\":null,\"file\":{\"filename\":\"\",\"max-size\":300,\"max-days\":0,\"max-backups\":0},\"enable-slow-log\":true,\"slow-query-file\":\"tidb-slow.log\",\"slow-threshold\":300,\"expensive-threshold\":10000,\"query-log-max-len\":4096,\"record-plan-in-slow-log\":1},\"security\":{\"skip-grant-table\":false,\"ssl-ca\":\"\",\"ssl-cert\":\"\",\"ssl-key\":\"\",\"require-secure-transport\":false,\"cluster-ssl-ca\":\"\",\"cluster-ssl-cert\":\"\",\"cluster-ssl-key\":\"\",\"cluster-verify-cn\":null,\"spilled-file-encryption-method\":\"plaintext\",\"enable-sem\":false,\"auto-tls\":false},\"status\":{\"status-host\":\"0.0.0.0\",\"metrics-addr\":\"\",\"status-port\":10080,\"metrics-interval\":15,\"report-status\":true,\"record-db-qps\":false},\"performance\":{\"max-procs\":0,\"max-memory\":0,\"server-memory-quota\":0,\"memory-usage-alarm-ratio\":0.8,\"stats-lease\":\"3s\",\"stmt-count-limit\":5000,\"feedback-probability\":0,\"query-feedback-limit\":512,\"pseudo-estimate-ratio\":0.8,\"force-priority\":\"NO_PRIORITY\",\"bind-info-lease\":\"3s\",\"txn-entry-size-limit\":6291456,\"txn-total-size-limit\":104857600,\"tcp-keep-alive\":true,\"tcp-no-delay\":true,\"cross-join\":true,\"run-auto-analyze\":true,\"distinct-agg-push-down\":false,\"committer-concurrency\":128,\"max-txn-ttl\":3600000,\"mem-profile-interval\":\"1m\",\"index-usage-sync-lease\":\"0s\",\"gogc\":100,\"enforce-mpp\":false},\"prepared-plan-cache\":{\"enabled\":false,\"capacity\":100,\"memory-guard-ratio\":0.1},\"opentracing\":{\"enable\":false,\"rpc-metrics\":false,\"sampler\":{\"type\":\"const\",\"param\":1,\"sampling-server-url\":\"\",\"max-operations\":0,\"sampling-refresh-interval\":0},\"reporter\":{\"queue-size\":0,\"buffer-flush-interval\":0,\"log-spans\":false,\"local-agent-host-port\":\"\"}},\"proxy-protocol\":{\"networks\":\"\",\"header-timeout\":5},\"pd-client\":{\"pd-server-timeout\":3},\"tikv-client\":{\"grpc-connection-count\":4,\"grpc-keepalive-time\":10,\"grpc-keepalive-timeout\":3,\"grpc-compression-type\":\"none\",\"commit-timeout\":\"41s\",\"async-commit\":{\"keys-limit\":256,\"total-key-size-limit\":4096,\"safe-window\":2000000000,\"allowed-clock-drift\":500000000},\"max-batch-size\":128,\"overload-threshold\":200,\"max-batch-wait-time\":0,\"batch-wait-size\":8,\"enable-chunk-rpc\":true,\"region-cache-ttl\":600,\"store-limit\":0,\"store-liveness-timeout\":\"1s\",\"copr-cache\":{\"capacity-mb\":1000},\"ttl-refreshed-txn-size\":33554432},\"binlog\":{\"enable\":false,\"ignore-error\":false,\"write-timeout\":\"15s\",\"binlog-socket\":\"\",\"strategy\":\"range\"},\"compatible-kill-query\":false,\"plugin\":{\"dir\":\"\",\"load\":\"\"},\"pessimistic-txn\":{\"max-retry-count\":256,\"deadlock-history-capacity\":10,\"deadlock-history-collect-retryable\":false},\"check-mb4-value-in-utf8\":true,\"max-index-length\":3072,\"index-limit\":64,\"table-column-count-limit\":1017,\"graceful-wait-before-shutdown\":0,\"alter-primary-key\":false,\"treat-old-version-utf8-as-utf8mb4\":true,\"enable-table-lock\":false,\"delay-clean-table-lock\":0,\"split-region-max-num\":1000,\"stmt-summary\":{\"enable\":true,\"enable-internal-query\":false,\"max-stmt-count\":3000,\"max-sql-length\":4096,\"refresh-interval\":1800,\"history-size\":24},\"repair-mode\":false,\"repair-table-list\":[],\"isolation-read\":{\"engines\":[\"tikv\",\"tiflash\",\"tidb\"]},\"max-server-connections\":0,\"new_collations_enabled_on_first_bootstrap\":false,\"experimental\":{\"allow-expression-index\":false},\"enable-collect-execution-info\":true,\"skip-register-to-dashboard\":false,\"enable-telemetry\":true,\"labels\":{},\"enable-global-index\":false,\"deprecate-integer-display-length\":false,\"enable-enum-length-limit\":true,\"stores-refresh-interval\":60,\"enable-tcp4-only\":false,\"enable-forwarding\":false}"]
	`
	s := NewStreamParser(strings.NewReader(logtxt))

	l, err := s.Next()
	assert.Nil(t, err)
	assert.NotNil(t, l)
	assert.Equal(t, s.Line, 2)
	assert.Equal(t, l.Message, "Mem Profile Tracker started")

	l, err = s.Next()
	assert.Nil(t, err)
	assert.NotNil(t, l)
	assert.Equal(t, s.Line, 3)
	assert.Equal(t, l.Message, "Welcome to TiDB.")

	l, err = s.Next()
	assert.NotNil(t, err)
	assert.Nil(t, l)
	assert.Equal(t, s.Line, 4)

	l, err = s.Next()
	assert.Nil(t, err)
	assert.NotNil(t, l)
	assert.Equal(t, s.Line, 5)
	assert.Equal(t, l.Message, "loaded config")
	assert.True(t, json.Valid([]byte(l.Fields[0].Value)))

	l, err = s.Next()
	assert.Nil(t, err)
	assert.Nil(t, l)
}
