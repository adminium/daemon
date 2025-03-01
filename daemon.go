package daemon

import (
	"github.com/adminium/logger"
	"sync/atomic"
)

var log = logger.NewLogger("daemon")

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
	err = d.init()
	if err != nil {
		return
	}
	for _, v := range d.conf.Processes {
		func(v *ProcessConfig) {
			p := newProcess(processArgs{
				name:            v.Name,
				command:         v.Command,
				restartInterval: v.getRestartInterval(d.conf),
				shell:           d.conf.getShell(),
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
	for _, v := range d.processes {
		if v.Cmd != nil {
			err := v.Cmd.Process.Kill()
			if err != nil {
				log.Errorf("kill process: %s error: %v", v.name, err)
			}
		}
		//	err := v.Cmd.Process.Signal(syscall.SIGTERM)
		//	if err != nil {
		//		log.Errorf("send SIGTERM to process: %s eror: %v", v.name, err)
		//		err = v.Cmd.Process.Kill()
		//		if err != nil {
		//			log.Errorf("kill process: %s error: %v", v.name, err)
		//		}
		//	} else {
		//		err = v.Cmd.Process.Release()
		//		if err != nil {
		//			log.Errorf("release process: %s error: %v", v.name, err)
		//		}
		//	}
		//}
	}
}
