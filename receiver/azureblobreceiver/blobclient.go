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

package azureblobreceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/azureblobreceiver"

import (
	"bytes"
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"go.uber.org/zap"
)

type BlobClient interface {
	ReadBlob(ctx context.Context, containerName string, blobName string) (*bytes.Buffer, error)
}

type AzureBlobClient struct {
	serviceClient *azblob.ServiceClient
	logger        *zap.Logger
}

func (bc *AzureBlobClient) getBlockBlob(containerName string, blobName string) (*azblob.BlockBlobClient, error) {
	containerClient, err := bc.serviceClient.NewContainerClient(containerName)
	if err != nil {
		return nil, err
	}

	return containerClient.NewBlockBlobClient(blobName)
}

func (bc *AzureBlobClient) ReadBlob(ctx context.Context, containerName string, blobName string) (*bytes.Buffer, error) {
	blockBlob, err := bc.getBlockBlob(containerName, blobName)
	if err != nil {
		return nil, err
	}
	defer func() {
		_, blobDeleteErr := blockBlob.Delete(ctx, nil)
		if blobDeleteErr != nil {
			bc.logger.Error(blobDeleteErr.Error())
		}
	}()

	get, err := blockBlob.Download(ctx, nil)
	if err != nil {
		return nil, err
	}

	downloadedData := &bytes.Buffer{}
	reader := get.Body(nil)
	defer reader.Close()

	_, err = downloadedData.ReadFrom(reader)

	return downloadedData, err
}

func NewBlobClient(connectionString string, logger *zap.Logger) (*AzureBlobClient, error) {
	serviceClient, err := azblob.NewServiceClientFromConnectionString(connectionString, nil)
	if err != nil {
		return nil, err
	}

	return &AzureBlobClient{
		serviceClient,
		logger,
	}, nil
}
