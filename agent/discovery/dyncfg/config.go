// SPDX-License-Identifier: GPL-3.0-or-later

package dyncfg

import (
	"github.com/khulnasoft/go.plugin/agent/confgroup"
	"github.com/khulnasoft/go.plugin/agent/functions"
	"github.com/khulnasoft/go.plugin/agent/module"
)

type Config struct {
	Plugin               string
	API                  KhulnasoftDyncfgAPI
	Functions            FunctionRegistry
	Modules              module.Registry
	ModuleConfigDefaults confgroup.Registry
}

type KhulnasoftDyncfgAPI interface {
	DynCfgEnable(string) error
	DynCfgReset() error
	DyncCfgRegisterModule(string) error
	DynCfgRegisterJob(_, _, _ string) error
	DynCfgReportJobStatus(_, _, _, _ string) error
	FunctionResultSuccess(_, _, _ string) error
	FunctionResultReject(_, _, _ string) error
}

type FunctionRegistry interface {
	Register(name string, reg func(functions.Function))
}

func validateConfig(cfg Config) error {
	return nil
}
