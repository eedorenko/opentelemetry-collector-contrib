// Copyright OpenTelemetry Authors
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

package azureblobexporter

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.uber.org/zap/zaptest"
)

var (
	testLogs = []byte(`{"resourceLogs":[{"resource":{"attributes":[{"key":"service.name","value":{"stringValue":"dotnet"}}]},"instrumentationLibraryLogs":[{"instrumentationLibrary":{},"logRecords":[{"timeUnixNano":"1643240673066096200","severityText":"Information","body":{"stringValue":"Message Body"},"flags":1,"traceId":"7b20d1349ef9b6d6f9d4d1d4a3ac2e82","spanId":"0c2ad924e1771630"}]}]}]}`)
)

// Test onLogData callback for the test logs data
func TestExporterLogDataCallback(t *testing.T) {
	logs := getTestLogs(t)

	blobClient := NewMockBlobClient()

	exporter := getLogsExporter(t, blobClient)

	assert.NoError(t, exporter.onLogData(context.Background(), logs))

	blobClient.AssertNumberOfCalls(t, "UploadData", 1)
}

func getLogsExporter(tb testing.TB, blobClient BlobClient) *logExporter {
	exporter := &logExporter{
		blobClient:    blobClient,
		logger:        zaptest.NewLogger(tb),
		logsMarshaler: plog.NewJSONMarshaler(),
	}

	return exporter
}

func getTestLogs(tb testing.TB) plog.Logs {
	logsUnmarshaler := plog.NewJSONUnmarshaler()
	logs, err := logsUnmarshaler.UnmarshalLogs(testLogs)
	require.NoError(tb, err, "Can't unmarshal testing logs data -> %s", err)
	return logs
}
