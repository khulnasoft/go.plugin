// SPDX-License-Identifier: GPL-3.0-or-later

package redis

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/blang/semver/v4"
	"regexp"
	"strings"
)

const precision = 1000 // float values multiplier and dimensions divisor

func (r *Redis) collect() (map[string]int64, error) {
	info, err := r.rdb.Info(context.Background(), "all").Result()
	if err != nil {
		return nil, err
	}

	if r.server == "" {
		s, v, err := extractServerVersion(info)
		if err != nil {
			return nil, fmt.Errorf("can not extract server app and version: %v", err)
		}
		r.server, r.version = s, v
		r.Debugf(`server="%s",version="%s"`, s, v)
	}

	if r.server != "redis" {
		return nil, fmt.Errorf("unsupported server app, want=redis, got=%s", r.server)
	}

	mx := make(map[string]int64)
	r.collectInfo(mx, info)
	r.collectPingLatency(mx)

	return mx, nil
}

// redis_version:6.0.9
var reVersion = regexp.MustCompile(`([a-z]+)_version:(\d+\.\d+\.\d+)`)

func extractServerVersion(info string) (string, *semver.Version, error) {
	var versionLine string
	for sc := bufio.NewScanner(strings.NewReader(info)); sc.Scan(); {
		line := sc.Text()
		if strings.Contains(line, "_version") {
			versionLine = strings.TrimSpace(line)
			break
		}
	}
	if versionLine == "" {
		return "", nil, errors.New("no version property")
	}

	match := reVersion.FindStringSubmatch(versionLine)
	if match == nil {
		return "", nil, fmt.Errorf("can not parse version property '%s'", versionLine)
	}

	server, version := match[1], match[2]
	ver, err := semver.New(version)
	if err != nil {
		return "", nil, err
	}

	return server, ver, nil
}
