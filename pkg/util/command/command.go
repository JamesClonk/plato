package command

import (
	"bytes"
	"io"
	"os"
	xc "os/exec"
	"strings"

	"github.com/JamesClonk/plato/pkg/util/color"
	"github.com/JamesClonk/plato/pkg/util/log"
	"github.com/JamesClonk/plato/pkg/util/terminal"
)

func Get(command []string) *xc.Cmd {
	return xc.Command(command[0], command[1:]...)
}

func GetInDir(command []string, dir string) *xc.Cmd {
	cmd := Get(command)
	cmd.Dir = dir
	return cmd
}

func ExecOutput(command []string) (string, error) {
	return execOutput(Get(command))
}

func RunOutput(command []string) string {
	return runOutput(Get(command))
}

func ExecOutputInDir(command []string, dir string) (string, error) {
	return execOutput(GetInDir(command, dir))
}

func RunOutputInDir(command []string, dir string) string {
	return runOutput(GetInDir(command, dir))
}

func execOutput(cmd *xc.Cmd) (string, error) {
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Debugf("Failed command: %v", strings.Join(cmd.Args, " "))
		return string(output), err
	}
	return string(output), nil
}

func runOutput(cmd *xc.Cmd) string {
	output, err := execOutput(cmd)
	if err != nil {
		log.Debugf("Failed command: %v", strings.Join(cmd.Args, " "))
		log.Errorf(color.Red("%s", output))
		log.Fatalf(color.Red("%v", err))
	}
	return output
}

func Exec(command []string) error {
	return exec(Get(command))
}

func Run(command []string) {
	run(Get(command))
}

func ExecInDir(command []string, dir string) error {
	return exec(GetInDir(command, dir))
}

func RunInDir(command []string, dir string) {
	run(GetInDir(command, dir))
}

func exec(cmd *xc.Cmd) error {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	if err != nil {
		log.Debugf("Failed command: %v", strings.Join(cmd.Args, " "))
	}
	return err
}

func run(cmd *xc.Cmd) {
	if err := exec(cmd); err != nil {
		log.Fatalf("Failed command: %s", color.Red("%v", err))
	}
}

func ExecWithLogfiles(command []string, logfile, errfile string) error {
	return execWithLogfiles(Get(command), logfile, errfile)
}

func RunWithLogfiles(command []string, logfile, errfile string) {
	runWithLogfiles(Get(command), logfile, errfile)
}

func execWithLogfiles(cmd *xc.Cmd, logfile, errfile string) error {
	logout, err := os.OpenFile(logfile, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0640)
	if err != nil {
		return err
	}
	defer logout.Close()
	logerr, err := os.OpenFile(errfile, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0640)
	if err != nil {
		return err
	}
	defer logerr.Close()

	cmd.Stdout = io.MultiWriter(terminal.NewWriter(logout), os.Stdout)
	cmd.Stderr = io.MultiWriter(terminal.NewWriter(logerr), os.Stderr)
	cmd.Stdin = os.Stdin

	err = cmd.Run()
	if err != nil {
		log.Debugf("Failed command: %v", strings.Join(cmd.Args, " "))
	}
	return err
}

func runWithLogfiles(cmd *xc.Cmd, logfile, errfile string) {
	if err := execWithLogfiles(cmd, logfile, errfile); err != nil {
		log.Debugf("Failed command: %v", strings.Join(cmd.Args, " "))
		log.Fatalf(color.Red("%v", err))
	}
}

func RunPiped(commands ...*xc.Cmd) string {
	if len(commands) < 1 {
		log.Fatalf(color.Fail("Not enough commands passed to RunPiped()"))
	}
	var stdout bytes.Buffer

	for c := range commands[:len(commands)-1] {
		if pipe, err := commands[c].StdoutPipe(); err != nil {
			log.Debugf("Failed command: %v", strings.Join(commands[c].Args, " "))
			log.Fatalf(color.Red("%v", err))
		} else {
			commands[c+1].Stdin = pipe
		}
	}
	commands[len(commands)-1].Stdout = &stdout

	for _, command := range commands {
		if err := command.Start(); err != nil {
			return stdout.String()
		}
	}

	for _, command := range commands {
		if err := command.Wait(); err != nil {
			return stdout.String()
		}
	}
	return stdout.String()
}
