package daemon

import "github.com/adminium/logger"

var log = logger.NewLogger("daemon")

func NewDaemon(conf *Config) *Daemon {
	return &Daemon{conf: conf}
}

type Daemon struct {
	conf      *Config
	processes []*Process
}

func (d *Daemon) init() (err error) {
	err = d.conf.parse()
	if err != nil {
		return
	}
	return
}

func (d *Daemon) Run() (err error) {

	err = d.init()
	if err != nil {
		return
	}

	for _, v := range d.conf.Processes {
		p :=
		return
	}

	return
}
