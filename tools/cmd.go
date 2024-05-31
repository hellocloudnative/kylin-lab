package tools

import (
	"bufio"
	"fmt"
	log "github.com/wonderivan/logger"
	"io"
	"os/exec"
	"sync"
	"syscall"
)

type Command struct {
	CmdStr          string
	Pid             int
	ExitCode        int
	StdOutput       string
	ErrOutput       string
	isPrintRealTime bool
}

const (
	ExitCodeDefault    = -999
	ExitCodeIoError    = -998
	ExitCodeStartError = -997
)

/*
*
创建命令执行实例
*/
func NewCmd(cmdStr string) *Command {
	return &Command{
		CmdStr:          cmdStr,
		ExitCode:        ExitCodeDefault,
		isPrintRealTime: false,
	}
}

/*
*
创建命令执行实例，并且实时打印输出
*/
func NewCmdWithPrint(cmdStr string) *Command {
	return &Command{
		CmdStr:          cmdStr,
		ExitCode:        ExitCodeDefault,
		isPrintRealTime: true,
	}
}

/*
*
执行命令
*/
func (cmd *Command) Start() {
	wg := &sync.WaitGroup{}
	wg.Add(2)
	cmdExec := exec.Command("sh", "-c", cmd.CmdStr)
	cmdExec.SysProcAttr = &syscall.SysProcAttr{Setpgid: true} // 将PGID设置成与PID相同的值 保证kill能够杀死孙子进程
	stdout, err := cmdExec.StdoutPipe()
	if err != nil {
		log.Error(fmt.Sprintf("get stdout failed: %s", err.Error()))
		cmd.ExitCode = ExitCodeIoError
		cmd.ErrOutput += err.Error()
		return
	}
	errout, err := cmdExec.StderrPipe()
	if err != nil {
		log.Error(fmt.Sprintf("get stdout failed: %s", err.Error()))
		cmd.ExitCode = ExitCodeIoError
		cmd.ErrOutput += err.Error()
		return
	}

	// 异步读取输出流
	go func() {
		defer wg.Done()
		cmd.readInputStream(stdout, false)
	}()
	go func() {
		defer wg.Done()
		cmd.readInputStream(errout, true)
	}()

	log.Info(fmt.Sprintf("start process %s", cmd.CmdStr))
	err = cmdExec.Start()
	if err != nil {
		log.Error(fmt.Sprintf("start process failed: %s", err.Error()))
		cmd.ExitCode = ExitCodeStartError
		cmd.ErrOutput += err.Error()
		return
	}
	cmd.Pid = cmdExec.Process.Pid
	wg.Wait()
	_ = cmdExec.Wait()
	cmd.ExitCode = cmdExec.ProcessState.ExitCode()
}

func (cmd *Command) Kill() error {
	log.Info(fmt.Sprintf("kill process %s,pid %d", cmd.CmdStr, cmd.Pid))
	return syscall.Kill(-cmd.Pid, syscall.SIGKILL)
}

/*
*
读取输出流，如果是错误输出流，默认实时打印
*/
func (cmd *Command) readInputStream(in io.ReadCloser, isError bool) {
	reader := bufio.NewReader(in)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if io.EOF != err {
				log.Error(err.Error())
			}
			break
		}
		if isError {
			cmd.ErrOutput += line
			log.Error(line)
		} else {
			cmd.StdOutput += line
			if cmd.isPrintRealTime {
				log.Info(line)
			}
		}
	}
}
