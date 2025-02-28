package daemon

import "github.com/adminium/logger"

var log = logger.NewLogger("daemon")

func NewDaemon(conf Config) *Daemon {
	return &Daemon{
		conf: conf,
		stop: make(chan struct{}),
	}
}

type Daemon struct {
	conf      Config
	processes []*Process
	stop      chan struct{}
}

func (d *Daemon) init() (err error) {
	err = d.conf.parse()
	if err != nil {
		return
	}
	return
}

func (d *Daemon) MakeLogDir() (err error) {
	return
}

func (d *Daemon) Run() (err error) {

	err = d.init()
	if err != nil {
		return
	}
	for _, v := range d.conf.Processes {
		func(v *ProcessConfig) {
			p := newProcess(processArgs{
				logDir:          d.conf.LogDir,
				name:            v.Name,
				command:         v.Command,
				restartInterval: v.getRestartInterval(d.conf),
				shell:           d.conf.getShell(),
			})
			go p.Run()
		}(v)
	}
	return
}

func (d *Daemon) Stop() {
	d.stop <- struct{}{}
	close(d.stop)
}
