// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/khulnasoft/go.plugin/agent"
	"github.com/khulnasoft/go.plugin/agent/executable"
	"github.com/khulnasoft/go.plugin/cli"
	"github.com/khulnasoft/go.plugin/logger"
	"github.com/khulnasoft/go.plugin/pkg/multipath"

	"github.com/jessevdk/go-flags"
	"golang.org/x/net/http/httpproxy"

	_ "github.com/khulnasoft/go.plugin/modules"
)

var (
	cd, _       = os.Getwd()
	name        = "go"
	userDir     = os.Getenv("KHULNASOFT_USER_CONFIG_DIR")
	stockDir    = os.Getenv("KHULNASOFT_STOCK_CONFIG_DIR")
	varLibDir   = os.Getenv("KHULNASOFT_LIB_DIR")
	lockDir     = os.Getenv("KHULNASOFT_LOCK_DIR")
	watchPath   = os.Getenv("KHULNASOFT_PLUGINS_GOD_WATCH_PATH")
	envLogLevel = os.Getenv("KHULNASOFT_LOG_LEVEL")

	version = "unknown"
)

func confDir(opts *cli.Option) multipath.MultiPath {
	if len(opts.ConfDir) > 0 {
		return opts.ConfDir
	}
	if userDir != "" || stockDir != "" {
		return multipath.New(
			userDir,
			stockDir,
		)
	}
	if executable.Directory != "" {
		return multipath.New(
			filepath.Join(executable.Directory, "/../../../../etc/khulnasoft"),
			filepath.Join(executable.Directory, "/../../../../usr/lib/khulnasoft/conf.d"),
		)
	}
	return multipath.New(
		filepath.Join(cd, "/../../../../etc/khulnasoft"),
		filepath.Join(cd, "/../../../../usr/lib/khulnasoft/conf.d"),
	)
}

func modulesConfDir(opts *cli.Option) (mpath multipath.MultiPath) {
	if len(opts.ConfDir) > 0 {
		return opts.ConfDir
	}
	if userDir != "" || stockDir != "" {
		if userDir != "" {
			mpath = append(mpath, filepath.Join(userDir, name))
		}
		if stockDir != "" {
			mpath = append(mpath, filepath.Join(stockDir, name))
		}
		return multipath.New(mpath...)
	}
	if executable.Directory != "" {
		return multipath.New(
			filepath.Join(executable.Directory, "/../../../../etc/khulnasoft", name),
			filepath.Join(executable.Directory, "/../../../../usr/lib/khulnasoft/conf.d", name),
		)
	}
	return multipath.New(
		filepath.Join(cd, "/../../../../etc/khulnasoft", name),
		filepath.Join(cd, "/../../../../usr/lib/khulnasoft/conf.d", name),
	)
}

func watchPaths(opts *cli.Option) []string {
	if watchPath == "" {
		return opts.WatchPath
	}
	return append(opts.WatchPath, watchPath)
}

func stateFile() string {
	if varLibDir == "" {
		return ""
	}
	return filepath.Join(varLibDir, "god-jobs-statuses.json")
}

func init() {
	// https://github.com/khulnasoft/khulnasoft/issues/8949#issuecomment-638294959
	if v := os.Getenv("TZ"); strings.HasPrefix(v, ":") {
		_ = os.Unsetenv("TZ")
	}
}

func main() {
	opts := parseCLI()

	if opts.Version {
		fmt.Printf("go.plugin, version: %s\n", version)
		return
	}

	if envLogLevel != "" {
		logger.Level.SetByName(envLogLevel)
	}

	if opts.Debug {
		logger.Level.Set(slog.LevelDebug)
	}

	a := agent.New(agent.Config{
		Name:              name,
		ConfDir:           confDir(opts),
		ModulesConfDir:    modulesConfDir(opts),
		ModulesSDConfPath: watchPaths(opts),
		VnodesConfDir:     confDir(opts),
		StateFile:         stateFile(),
		LockDir:           lockDir,
		RunModule:         opts.Module,
		MinUpdateEvery:    opts.UpdateEvery,
	})

	a.Debugf("plugin: name=%s, version=%s", a.Name, version)
	if u, err := user.Current(); err == nil {
		a.Debugf("current user: name=%s, uid=%s", u.Username, u.Uid)
	}

	cfg := httpproxy.FromEnvironment()
	a.Infof("env HTTP_PROXY '%s', HTTPS_PROXY '%s'", cfg.HTTPProxy, cfg.HTTPSProxy)

	a.Run()
}

func parseCLI() *cli.Option {
	opt, err := cli.Parse(os.Args)
	if err != nil {
		var flagsErr *flags.Error
		if errors.As(err, &flagsErr) && errors.Is(flagsErr.Type, flags.ErrHelp) {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}
	return opt
}
