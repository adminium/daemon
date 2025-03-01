package daemon

import (
	"errors"
	"fmt"
	"github.com/adminium/logger"
	"os"
	"os/exec"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

type processArgs struct {
	name            string
	command         string
	shell           string
	restartInterval time.Duration
	stopSignal      os.Signal
	stopWaitTime    time.Duration
	stopped         *atomic.Bool
}

func newProcess(args processArgs) *Process {
	return &Process{
		processArgs: args,
		log:         logger.NewLogger(fmt.Sprintf("process::%s", args.name)),
	}
}

type Process struct {
	processArgs
	*exec.Cmd
	log *logger.Logger
}

func (p *Process) Run() {
Start:
	p.log.Infof("start")
	p.Cmd = exec.Command(p.shell, "-c", p.command)
	p.Cmd.Env = os.Environ()
	p.Cmd.Stdout = os.Stdout
	p.Cmd.Stderr = os.Stderr
	err := p.Cmd.Start()
	if err != nil {
		p.log.Errorf("exec process error: %s, restart after: %s", err, p.restartInterval)
		time.Sleep(p.restartInterval)
		goto Start
	}

	_, err = p.Cmd.Process.Wait()
	if err != nil {
		p.log.Errorf("wait process error: %s", err)
	}
	err = p.Cmd.Process.Kill()
	if err == nil {
		p.log.Errorf("kill process error: %s", err)
	}

	if p.stopped.Load() {
		return
	}

	p.log.Infof("restrt process after: %s", p.restartInterval)
	time.Sleep(p.restartInterval)

	goto Start
}

func (p *Process) Stop(wg *sync.WaitGroup) {
	defer wg.Done()
	go func() {
		err := p.Cmd.Process.Signal(p.stopSignal)
		if err != nil {
			log.Errorf("send signal: %s to process: %s eror: %v", p.stopSignal.String(), p.name, err)
			return
		}
	}()

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	now := time.Now()
	for {
		select {
		case <-ticker.C:
			err := p.Cmd.Process.Signal(syscall.Signal(0))
			if err != nil {
				if errors.Is(err, os.ErrProcessDone) {
					return
				}
			}
			if time.Since(now) > p.stopWaitTime {
				err = p.Cmd.Process.Kill()
				if err != nil {
					p.log.Warnf("kill process: %s eror: %v", p.name, err)
				}
				return
			}
		}
	}
}
