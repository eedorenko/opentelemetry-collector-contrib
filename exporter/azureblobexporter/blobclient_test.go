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
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/config"
	"go.uber.org/zap/zaptest"
)

const (
	goodConnectionString = "DefaultEndpointsProtocol=https;AccountName=accountName;AccountKey=+idLkHYcL0MUWIKYHm2j4Q==;EndpointSuffix=core.windows.net"
	badConnectionString  = "DefaultEndpointsProtocol=https;AccountName=accountName;AccountKey=accountkey;EndpointSuffix=core.windows.net"
)

func TestNewBlobClient(t *testing.T) {
	blobClient, err := NewBlobClient(goodConnectionString, logsContainerName, zaptest.NewLogger(t))

	require.Nil(t, err)
	require.NotNil(t, blobClient)
	assert.NotNil(t, blobClient.containerClient)
}

func TestNewBlobClientError(t *testing.T) {
	blobClient, err := NewBlobClient(badConnectionString, logsContainerName, zaptest.NewLogger(t))

	assert.NotNil(t, err)
	assert.Nil(t, blobClient)
}

func TestGenerateBlobName(t *testing.T) {
	blobClient, err := NewBlobClient(goodConnectionString, logsContainerName, zaptest.NewLogger(t))
	require.Nil(t, err)

	blobName := blobClient.generateBlobName(config.LogsDataType)
	assert.True(t, strings.Contains(blobName, fmt.Sprintf("%s-", config.LogsDataType)))
}

func TestCheckOrCreateContainer(t *testing.T) {
	blobClient, err := NewBlobClient(goodConnectionString, logsContainerName, zaptest.NewLogger(t))
	require.Nil(t, err)

	err = blobClient.checkOrCreateContainer()

	assert.NotNil(t, err)

	assert.False(t, strings.Contains(err.Error(), containerNotFoundError))

}
