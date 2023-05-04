/*
   Copyright The containerd Authors.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package images

import (
	"context"
	"io"

	"github.com/containerd/containerd/v2/content"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

// Importer is the interface for image importer.
type Importer interface {
	// Import imports an image from a tar stream.
	Import(ctx context.Context, store content.Store, reader io.Reader) (ocispec.Descriptor, error)
}

// Exporter is the interface for image exporter.
type Exporter interface {
	// Export exports an image to a tar stream.
	Export(ctx context.Context, store content.Provider, desc ocispec.Descriptor, writer io.Writer) error
}
