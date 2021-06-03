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
func launcher(wg *sync.WaitGroup, exitChan chan bool, shell string, scriptPath string,
	launchLogger *log.Logger, scriptLogger *log.Logger) {

	defer wg.Done()
	defer common.SignalFinished(exitChan)

	if _, err := os.Stat(scriptPath); err != nil {
		if os.IsNotExist(err) {
			launchLogger.Printf("%s does not exist", scriptPath)
			return
		}
	}

	outputFile, err := os.Create("output.txt")
	if err != nil {
		launchLogger.Println(err.Error())
		return
	}
	defer outputFile.Close()
	errFile, err := os.Create("err.txt")
	if err != nil {
		launchLogger.Println(err.Error())
		return
	}
	defer errFile.Close()

	var scriptCmd *exec.Cmd
	switch shell {
	case CMD.String():
		scriptCmd = exec.Command("cmd.exe", "/C", scriptPath)
	case POWERSHELL.String():
		shellCommand := fmt.Sprintf("[Console]::OutputEncoding = [Text.Encoding]::UTF8; .\\%s", scriptPath)
		scriptCmd = exec.Command("powershell.exe", "-NoProfile", "-NonInteractive", "-ExecutionPolicy", "Bypass", "-Command", shellCommand)
	case PWSH.String():
	default:
		launchLogger.Println("shell type not supported")
		return
	}
	scriptCmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: windows.CREATE_NEW_PROCESS_GROUP}
	// Note: Go writes the output in utf8 WITHOUT a bom. No need for any encoding conversions
	scriptCmd.Stdout = outputFile
	scriptCmd.Stderr = errFile

	for i := 0; i < len(scriptCmd.Args); i++ {
		launchLogger.Printf("args %v: %v\n", i, scriptCmd.Args[i])
	}
	err = scriptCmd.Run()
	if err != nil {
		launchLogger.Printf("cmd.Run() failed with %s\n", err)
	}
	launchLogger.Println("command finished")

	resultVal := scriptCmd.ProcessState.ExitCode()
	launchLogger.Printf("script exit code: %v\n", resultVal)
	common.ExitLauncher(resultVal, "result.txt", launchLogger)
}

func main() {
	var logPath, shell, scriptPath string
	var debug, daemon bool
	const logFlag = "log"
	const shellFlag = "shell"
	const scriptPathFlag = "script"
	const debugFlag = "debug"
	const daemonFlag = "daemon"
	flag.StringVar(&logPath, logFlag, "", "full path of the log file")
	flag.StringVar(&shell, shellFlag, "cmd", "Windows shell type")
	flag.StringVar(&scriptPath, scriptPathFlag, "", "full path of the script to be launched")
	flag.BoolVar(&debug, debugFlag, false, "noisy output to log")
	flag.BoolVar(&daemon, daemonFlag, false, "Free binary from parent process")
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
	logFile, logErr := os.Create(logPath)
	if logErr != nil {
		fmt.Fprintf(os.Stderr, "Unable to create log file: %s", logErr)
		return
	}
	defer logFile.Close()
	mainLogger, _, launchLogger, scriptLogger := common.PrepareLogging(logFile, debug)

	for key, val := range defined {
		mainLogger.Printf("%v: %v", key, val)
	}
	mainLogger.Printf("Main pid is: %v\n", os.Getpid())
	mainLogger.Printf("Parent pid is: %v\n", os.Getppid())

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go common.SignalCatcher(sigChan, mainLogger)

	exitChan := make(chan bool)
	var wg sync.WaitGroup
	wg.Add(1)
	go launcher(&wg, exitChan, shell, scriptPath, launchLogger, scriptLogger)
	// TEMP until we add heartbeat: Must access the channel else it blocks launcher
	channelResult := <-exitChan
	mainLogger.Printf("exit chan is %v\n", channelResult)
	wg.Wait()
	signal.Stop(sigChan)
	close(sigChan)
	mainLogger.Println("done.")
}
