package config

import (
	"io"

	srvconfig "github.com/containerd/containerd/services/server/config"
	"github.com/pelletier/go-toml"
)

// File is a wrapper of server config for printing out.
type File struct {
	*srvconfig.Config
	// Plugins overrides `Plugins map[string]toml.Tree` in server config.
	Plugins map[string]interface{} `toml:"plugins"`
}

// WriteTo marshals the config to the provided writer
func (c *File) WriteTo(w io.Writer) (int64, error) {
	return 0, toml.NewEncoder(w).Encode(c)
}
