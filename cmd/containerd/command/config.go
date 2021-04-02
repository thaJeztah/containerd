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

package command

import (
	gocontext "context"
	"os"

	"github.com/containerd/containerd/config"
	"github.com/containerd/containerd/pkg/timeout"
	"github.com/containerd/containerd/services/server"
	srvconfig "github.com/containerd/containerd/services/server/config"
	"github.com/urfave/cli"
)

func outputConfig(cfg *srvconfig.Config) error {
	out := &config.File{
		Config: cfg,
	}

	plugins, err := server.LoadPlugins(gocontext.Background(), out.Config)
	if err != nil {
		return err
	}
	if len(plugins) != 0 {
		out.Plugins = make(map[string]interface{})
		for _, p := range plugins {
			if p.Config == nil {
				continue
			}

			pc, err := out.Decode(p)
			if err != nil {
				return err
			}

			out.Plugins[p.URI()] = pc
		}
	}

	if out.Timeouts == nil {
		out.Timeouts = make(map[string]string)
	}
	timeouts := timeout.All()
	for k, v := range timeouts {
		if out.Timeouts[k] == "" {
			out.Timeouts[k] = v.String()
		}
	}

	// for the time being, keep the defaultConfig's version set at 1 so that
	// when a config without a version is loaded from disk and has no version
	// set, we assume it's a v1 config.  But when generating new configs via
	// this command, generate the v2 config
	out.Config.Version = 2

	// remove overridden Plugins type to avoid duplication in output
	out.Config.Plugins = nil

	_, err = out.WriteTo(os.Stdout)
	return err
}

var configCommand = cli.Command{
	Name:  "config",
	Usage: "information on the containerd config",
	Subcommands: []cli.Command{
		{
			Name:  "default",
			Usage: "see the output of the default config",
			Action: func(context *cli.Context) error {
				return outputConfig(config.Default())
			},
		},
		{
			Name:  "dump",
			Usage: "see the output of the final main config with imported in subconfig files",
			Action: func(context *cli.Context) error {
				config := config.Default()
				if err := srvconfig.LoadConfig(context.GlobalString("config"), config); err != nil && !os.IsNotExist(err) {
					return err
				}

				return outputConfig(config)
			},
		},
	},
}
