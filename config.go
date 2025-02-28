package daemon

import (
	"fmt"
	"strings"
	"time"
)

var (
	DefaultRestartInterval = 1 * time.Second
)

type Config struct {
	LogDir          string           `toml:"log_dir"`
	Shell           string           `toml:"shell"`
	RestartInterval string           `toml:"restart_interval"`
	Processes       []*ProcessConfig `toml:"processes"`
	restartInterval time.Duration
}

func (c *Config) getShell() string {
	if c.Shell != "" {
		return c.Shell
	}
	return "sh"
}

func (c *Config) parse() (err error) {

	if c.LogDir == "" {
		err = fmt.Errorf("log dir is required")
		return
	}

	if c.RestartInterval != "" {
		c.restartInterval, err = time.ParseDuration(c.RestartInterval)
		if err != nil {
			err = fmt.Errorf("parse restart interval: %s error: %s", c.RestartInterval, err)
			return
		}
	}

	m := map[string]struct{}{}
	for _, v := range c.Processes {
		v.Name = strings.TrimSpace(v.Name)
		if v.Name == "" {
			err = fmt.Errorf("process name is required")
			return
		}
		v.Command = strings.TrimSpace(v.Command)
		if v.Command == "" {
			err = fmt.Errorf("process: %s command is empty", v.Name)
			return
		}
		if _, ok := m[v.Name]; ok {
			err = fmt.Errorf("process: %s is duplicated", v.Name)
			return
		}
		m[v.Name] = struct{}{}
		err = v.parse()
		if err != nil {
			err = fmt.Errorf("parse process: %s config failed: %s", v.Name, err)
			return
		}
	}

	return
}

type ProcessConfig struct {
	Name            string `toml:"name"`
	Command         string `toml:"command"`
	RestartInterval string `toml:"restart_interval"`
	restartInterval time.Duration
}

func (p *ProcessConfig) parse() (err error) {
	if p.RestartInterval != "" {
		p.restartInterval, err = time.ParseDuration(p.RestartInterval)
		if err != nil {
			return
		}
	}
	return
}

func (p *ProcessConfig) getRestartInterval(conf Config) time.Duration {
	if p.restartInterval > 0 {
		return p.restartInterval
	}
	if conf.restartInterval > 0 {
		return conf.restartInterval
	}
	return DefaultRestartInterval
}
