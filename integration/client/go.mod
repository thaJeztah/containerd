module github.com/containerd/containerd/integration/client

go 1.15

require (
	github.com/Microsoft/hcsshim v0.8.16
	github.com/Microsoft/hcsshim/test v0.0.0-20210408205431-da33ecd607e1
	github.com/containerd/cgroups v0.0.0-20210114181951-8a68de567b68
	// the actual version of containerd is replaced with the code at the root of this repository
	github.com/containerd/containerd v1.5.0-beta.4
	github.com/containerd/go-runc v0.0.0-20201020171139-16b287bc67d0
	github.com/containerd/ttrpc v1.0.2
	github.com/containerd/typeurl v1.0.1
	github.com/gogo/protobuf v1.3.2
	github.com/opencontainers/go-digest v1.0.0
	github.com/opencontainers/image-spec v1.0.1
	github.com/opencontainers/runtime-spec v1.0.3-0.20200929063507-e6143ca7d51d
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.7.0
	golang.org/x/sys v0.0.0-20210324051608-47abb6519492
	gotest.tools/v3 v3.0.3
)

replace (
	// use the containerd module from this repository instead of downloading
	//
	// IMPORTANT: this replace rule ONLY replaces containerd itself; dependencies
	// in the "require" section above are still taken into account for version
	// resolution if newer.
	github.com/containerd/containerd => ../../

	// Ignore the rules below, they're to simplify the reproducer
	github.com/prometheus/procfs => github.com/prometheus/procfs v0.2.0
	gopkg.in/yaml.v2 => gopkg.in/yaml.v2 v2.3.0
)
