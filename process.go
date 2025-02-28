package daemon

import (
	"fmt"
	"github.com/adminium/logger"
	"os"
	"os/exec"
	"sync/atomic"
	"time"
)

type processArgs struct {
	name            string
	command         string
	restartInterval time.Duration
	shell           string
	stopped         *atomic.Bool
}

type Process struct {
	processArgs
	*exec.Cmd
	log *logger.Logger
}

func newProcess(args processArgs) *Process {
	return &Process{
		processArgs: args,
		log:         logger.NewLogger(fmt.Sprintf("process::%s", args.name)),
	}
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
