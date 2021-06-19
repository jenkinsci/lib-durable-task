// +build windows

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"golang.org/x/sys/windows"
	"jenkinsci.org/plugins/durabletask/common"
)

// Launches the script in a new session and waits for its completion.
func launcher(wg *sync.WaitGroup, exitChan chan bool,
	executable string, args string, resultPath string, outputPath string,
	launchLogger *log.Logger, scriptLogger *log.Logger) {

	defer wg.Done()
	defer common.SignalFinished(exitChan)

	var scriptCmd *exec.Cmd
	var sysAttr windows.SysProcAttr

	sysAttr.CreationFlags = windows.CREATE_NEW_PROCESS_GROUP
	if executable == "cmd" {
		scriptCmd = exec.Command(executable)
		sysAttr.CmdLine = executable + " " + args
	} else {
		// Delimiter is a common WITH a space to handle complex arguments that require commas. For example, an argument that includes a method
		// call with args separated by a comma. The method args can be separated by commas with NO space so they do not accidentally get split
		scriptCmd = exec.Command(executable, strings.Split(args, ", ")...)
	}
	scriptCmd.SysProcAttr = &sysAttr

	if outputPath != "" {
		// capturing output
		outputFile, err := os.Create(outputPath)
		if err != nil {
			launchLogger.Println(err.Error())
			common.RecordExit(-2, resultPath, launchLogger)
			return
		}
		defer outputFile.Close()

		// Note: Go writes the output in utf8 WITHOUT a bom. No need for any encoding conversions
		scriptCmd.Stdout = outputFile
		scriptCmd.Stderr = scriptLogger.Writer()
	} else {
		scriptCmd.Stdout = scriptLogger.Writer()
		scriptCmd.Stderr = scriptCmd.Stdout
	}

	for i := 0; i < len(scriptCmd.Args); i++ {
		launchLogger.Printf("args %v: %v\n", i, scriptCmd.Args[i])
	}
	err := scriptCmd.Start()
	if common.CheckIfErr(scriptLogger, err) {
		common.RecordExit(-2, resultPath, scriptLogger)
		return
	}
	pid := scriptCmd.Process.Pid
	launchLogger.Printf("launched %v\n", pid)
	err = scriptCmd.Wait()
	common.CheckIfErr(scriptLogger, err)
	resultVal := scriptCmd.ProcessState.ExitCode()
	launchLogger.Printf("script exit code: %v\n", resultVal)

	common.RecordExit(resultVal, resultPath, launchLogger)
}

func main() {
	var controlDir, resultPath, logPath, executable, args, outputPath string
	var debug, daemon bool
	const controlFlag = "controldir"
	const resultFlag = "result"
	const logFlag = "log"
	const executableFlag = "executable"
	const argsFlag = "args"
	const outputFlag = "output"
	const debugFlag = "debug"
	const daemonFlag = "daemon"
	flag.StringVar(&controlDir, controlFlag, "", "working directory")
	flag.StringVar(&resultPath, resultFlag, "", "full path of the result file")
	flag.StringVar(&logPath, logFlag, "", "full path of the log file")
	flag.StringVar(&executable, executableFlag, "", "path to the executable being launched")
	flag.StringVar(&args, argsFlag, "", "(optional) argument(s) to the executable, separated by commas WITH spaces")
	flag.StringVar(&outputPath, outputFlag, "", "(optional) if recording output, full path of the output file")
	flag.BoolVar(&debug, debugFlag, false, "(optional) noisy output to log")
	flag.BoolVar(&daemon, daemonFlag, false, "(optional) Free binary from parent process")
	flag.Parse()

	// Validate that the required flags were all command-line defined
	required := []string{controlFlag, resultFlag, logFlag, executableFlag}
	defined := make(map[string]string)
	flag.Visit(func(f *flag.Flag) {
		defined[f.Name] = f.Value.String()
	})
	if !common.ValidateFlags(defined, required) {
		os.Exit(-2)
	}

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
		common.RecordExit(-2, resultPath, log.Default())
		return
	}
	defer logFile.Close()
	mainLogger, hbLogger, launchLogger, scriptLogger := common.PrepareLogging(logFile, debug)

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
	wg.Add(2)
	go launcher(&wg, exitChan, executable, args, resultPath, outputPath, launchLogger, scriptLogger)
	go common.Heartbeat(&wg, exitChan, controlDir, resultPath, logPath, hbLogger)
	mainLogger.Println("about to wait")
	wg.Wait()
	mainLogger.Println("done waiting")
	signal.Stop(sigChan)
	close(sigChan)
	mainLogger.Println("done.")
}
