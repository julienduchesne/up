// Copyright 2021 Upbound Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package controlplane

import (
	"github.com/alecthomas/kong"

	cp "github.com/upbound/up-sdk-go/service/controlplanes"

	"github.com/upbound/up/cmd/up/controlplane/kubeconfig"
	"github.com/upbound/up/cmd/up/controlplane/pkg"
	"github.com/upbound/up/cmd/up/controlplane/pullsecret"
	"github.com/upbound/up/internal/feature"
	"github.com/upbound/up/internal/upbound"
)

// BeforeReset is the first hook to run.
func (c *Cmd) BeforeReset(p *kong.Path, maturity feature.Maturity) error {
	return feature.HideMaturity(p, maturity)
}

// AfterApply constructs and binds a control plane client to any subcommands
// that have Run() methods that receive it.
func (c *Cmd) AfterApply(kongCtx *kong.Context) error {
	upCtx, err := upbound.NewFromFlags(c.Flags)
	if err != nil {
		return err
	}
	cfg, err := upCtx.BuildSDKConfig(upCtx.Profile.Session)
	if err != nil {
		return err
	}
	kongCtx.Bind(upCtx)
	kongCtx.Bind(cp.NewClient(cfg))
	return nil
}

// Cmd contains commands for interacting with control planes.
type Cmd struct {
	Create createCmd `cmd:"" maturity:"alpha" help:"Create a hosted control plane."`
	Delete deleteCmd `cmd:"" maturity:"alpha" help:"Delete a control plane."`
	List   listCmd   `cmd:"" maturity:"alpha" help:"List control planes for the account."`
	Get    getCmd    `cmd:"" maturity:"alpha" help:"Get a single control plane."`

	Configuration pkg.Cmd `cmd:"" set:"package_type=Configuration" help:"Manage Configurations."`
	Provider      pkg.Cmd `cmd:"" set:"package_type=Provider" help:"Manage Providers."`

	PullSecret pullsecret.Cmd `cmd:"" help:"Manage package pull secrets."`

	Kubeconfig kubeconfig.Cmd `cmd:"" maturity:"alpha" name:"kubeconfig" help:"Manage control plane kubeconfig data."`

	// Common Upbound API configuration
	Flags upbound.Flags `embed:""`
}
