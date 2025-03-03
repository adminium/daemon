package daemon

import (
	"github.com/adminium/exit"
	"github.com/adminium/logger"
	"sync"
	"sync/atomic"
)

var Log logger.EventLogger = logger.NewLogger("daemon")

var NewLogger = func(module string) logger.EventLogger {
	return logger.NewLogger(module)
}

func NewDaemon(conf Config) *Daemon {
	return &Daemon{
		conf:    conf,
		stopped: &atomic.Bool{},
	}
}

type Daemon struct {
	conf      Config
	processes []*Process
	stopped   *atomic.Bool
}

func (d *Daemon) init() (err error) {
	err = d.conf.parse()
	if err != nil {
		return
	}
	return
}

func (d *Daemon) Run() (err error) {
	err = d.Start()
	if err != nil {
		return
	}
	exit.Clean(func() {
		d.Stop()
	})
	exit.Wait()
	return
}

func (d *Daemon) Start() (err error) {
	err = d.init()
	if err != nil {
		return
	}
	for _, v := range d.conf.Processes {
		func(v *ProcessConfig) {
			p := newProcess(processArgs{
				name:            v.Name,
				command:         v.Command,
				shell:           d.conf.getShell(),
				restartInterval: v.getRestartInterval(d.conf),
				stopSignal:      v.getStopSignal(d.conf),
				stopWaitTime:    v.getStopWaitTime(d.conf),
				stopped:         d.stopped,
			})
			d.processes = append(d.processes, p)
			go p.Run()
		}(v)
	}
	return
}

func (d *Daemon) Stop() {
	if !d.stopped.CompareAndSwap(false, true) {
		return
	}
	wg := &sync.WaitGroup{}
	wg.Add(len(d.processes))
	for _, v := range d.processes {
		p := v
		go p.Stop(wg)
	}
	wg.Wait()
}
