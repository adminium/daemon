package daemon

import (
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"
)

var (
	DefaultRestartWaitTime = 3 * time.Second
	DefaultStopSignal      = syscall.SIGTERM
	DefaultStopWaitTime    = 10 * time.Second
)

type Config struct {
	Shell           string           `toml:"shell"`
	RestartWaitTime string           `toml:"restart_wait_time"`
	StopSignal      string           `toml:"stop_signal"`
	StopWaitTime    string           `toml:"stop_wait_time"`
	Processes       []*ProcessConfig `toml:"processes"`
	stopSignal      os.Signal
	stopWaitTime    time.Duration
	restartInterval time.Duration
}

func (c *Config) getShell() string {
	if c.Shell != "" {
		return c.Shell
	}
	return "sh"
}

func (c *Config) parse() (err error) {

	if c.RestartWaitTime != "" {
		c.restartInterval, err = time.ParseDuration(c.RestartWaitTime)
		if err != nil {
			err = fmt.Errorf("parse restart interval: %s error: %s", c.RestartWaitTime, err)
			return
		}
	}

	if c.StopSignal != "" {
		s, ok := signals[c.StopSignal]
		if !ok {
			err = fmt.Errorf("unknown stop signal: %s", c.StopSignal)
		}
		c.stopSignal = s
	}

	if c.StopWaitTime != "" {
		c.stopWaitTime, err = time.ParseDuration(c.StopWaitTime)
		if err != nil {
			err = fmt.Errorf("parse stop wait time: %s error: %s", c.StopWaitTime, err)
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
	RestartWaitTime string `toml:"restart_wait_time"`
	StopSignal      string `toml:"stop_signal"`
	StopWaitTime    string `toml:"stop_wait_time"`
	stopSignal      os.Signal
	stopWaitTime    time.Duration
	restartInterval time.Duration
}

func (p *ProcessConfig) parse() (err error) {
	if p.RestartWaitTime != "" {
		p.restartInterval, err = time.ParseDuration(p.RestartWaitTime)
		if err != nil {
			return
		}
	}
	if p.StopSignal != "" {
		s, ok := signals[p.StopSignal]
		if !ok {
			err = fmt.Errorf("unknown stop signal: %s", p.StopSignal)
		}
		p.stopSignal = s
	}

	if p.StopWaitTime != "" {
		p.stopWaitTime, err = time.ParseDuration(p.StopWaitTime)
		if err != nil {
			err = fmt.Errorf("parse stop wait time: %s error: %s", p.StopWaitTime, err)
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
	return DefaultRestartWaitTime
}

func (p *ProcessConfig) getStopSignal(conf Config) os.Signal {
	if p.stopSignal != nil {
		return p.stopSignal
	}
	if conf.stopSignal != nil {
		return conf.stopSignal
	}
	return DefaultStopSignal
}

func (p *ProcessConfig) getStopWaitTime(conf Config) time.Duration {
	if p.stopWaitTime > 0 {
		return p.stopWaitTime
	}
	if conf.stopWaitTime > 0 {
		return conf.stopWaitTime
	}
	return DefaultStopWaitTime
}
