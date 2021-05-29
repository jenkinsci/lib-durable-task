// +build windows

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"

	"golang.org/x/sys/windows"
	"jenkinsci.org/plugins/durabletask/common"
)

var logger *log.Logger

type Shell int

const (
	CMD Shell = iota
	POWERSHELL
	PWSH
)

func (shell Shell) String() string {
	switch shell {
	case CMD:
		return "cmd"
	case POWERSHELL:
		return "powershell"
	case PWSH:
		return "pwsh"
	default:
		return "UNKNOWN"
	}
}

// Launches the script in a new session and waits for its completion.
func launcher(wg *sync.WaitGroup, shell string, scriptPath string) {
	defer wg.Done()

	if _, err := os.Stat(scriptPath); err != nil {
		if os.IsNotExist(err) {
			logger.Printf("%s does not exist", scriptPath)
			return
		}
	}

	outputFile, err := os.Create("output.txt")
	if err != nil {
		logger.Println(err.Error())
		return
	}
	defer outputFile.Close()
	errFile, err := os.Create("err.txt")
	if err != nil {
		logger.Println(err.Error())
		return
	}
	defer errFile.Close()

	var scriptCmd *exec.Cmd
	switch shell {
	case CMD.String():
		scriptCmd = exec.Command("cmd.exe", "/C", scriptPath)
		scriptCmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: windows.CREATE_NEW_PROCESS_GROUP}
	case POWERSHELL.String():
		shellCommand := fmt.Sprintf("[Console]::OutputEncoding = [Text.Encoding]::UTF8; .\\%s", scriptPath)
		logger.Printf("powershell command: %s", shellCommand)
		scriptCmd = exec.Command("powershell.exe", "-NoProfile", "-NonInteractive", "-ExecutionPolicy", "Bypass", "-Command", shellCommand)
		scriptCmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: windows.CREATE_NEW_PROCESS_GROUP}
	case PWSH.String():
	default:
		logger.Println("shell type not supported")
		return
	}
	// Note: Go writes the output in utf8 WITHOUT a bom. No need for any encoding conversions
	scriptCmd.Stdout = outputFile
	scriptCmd.Stderr = errFile

	logger.Println("about to launch command")
	err = scriptCmd.Run()
	if err != nil {
		logger.Printf("cmd.Run() failed with %s\n", err)
	}
	logger.Println("command finished")

	resultVal := scriptCmd.ProcessState.ExitCode()
	logger.Printf("script exit code: %v\n", resultVal)
	common.ExitLauncher(resultVal, "result.txt", logger)
}

func main() {
	var daemon bool
	var shell, scriptPath string
	const daemonFlag = "daemon"
	const shellFlag = "shell"
	const scriptPathFlag = "path"
	flag.BoolVar(&daemon, daemonFlag, false, "Free binary from parent process")
	flag.StringVar(&shell, shellFlag, "cmd", "Windows shell type")
	flag.StringVar(&scriptPath, scriptPathFlag, "", "full path of the script to be launched")
	flag.Parse()

	// Validate that the required flags were all command-line defined
	required := []string{scriptPathFlag}
	defined := make(map[string]string)
	flag.Visit(func(f *flag.Flag) {
		defined[f.Name] = f.Value.String()
	})
	common.ValidateFlags(defined, required)

	fmt.Fprintf(os.Stdout, "Parent pid is: %v\n", os.Getppid())

	if daemon {
		fmt.Fprintf(os.Stdout, "1st launch pid is: %v\n", os.Getpid())
		rebuiltArgs := common.RebuildArgs(defined, daemonFlag)
		doubleLaunchCmd := exec.Command(os.Args[0], rebuiltArgs...)
		doubleLaunchCmd.Stdout = nil
		doubleLaunchCmd.Stderr = nil
		doubleLaunchCmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: windows.DETACHED_PROCESS | windows.CREATE_NEW_PROCESS_GROUP}
		doubleLaunchErr := doubleLaunchCmd.Start()
		if doubleLaunchErr != nil {
			panic("Double launch failed, exiting")
		}
		return
	}
	// Prepare logging
	logFile, logErr := os.Create("logging.txt")
	if logErr != nil {
		fmt.Fprintf(os.Stderr, "Unable to create log file: %s", logErr)
		return
	}
	defer logFile.Close()
	logger = log.New(logFile, "MAIN ", log.Lmicroseconds|log.Lshortfile)
	logger.Printf("binary pid is: %v\n", os.Getpid())
	logger.Printf("parent pid is: %v\n", os.Getppid())

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go common.SignalCatcher(sigChan, logger)

	var wg sync.WaitGroup
	wg.Add(1)
	go launcher(&wg, shell, scriptPath)
	wg.Wait()
	signal.Stop(sigChan)
	close(sigChan)
	logger.Println("done.")
}
