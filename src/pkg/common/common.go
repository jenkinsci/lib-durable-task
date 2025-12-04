/*
 * The MIT License
 *
 * Copyright 2021 CloudBees, Inc.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE.
 */

package common

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

// Convenience function to log errors
func CheckIfErr(logger *log.Logger, err error) bool {
	if err != nil {
		logger.Println(err.Error())
		return true
	}
	return false
}

func SignalFinished(exitChan chan bool) {
	exitChan <- true
}

// Validate that the required flags were all defined
func ValidateFlags(defined map[string]string, required []string) bool {
	var missing []string
	for _, reqFlag := range required {
		if _, exists := defined[reqFlag]; !exists {
			missing = append(missing, reqFlag)
		}
	}

	if len(missing) > 0 {
		fmt.Println("The following required flags are missing:")
		for _, missingFlag := range missing {
			fmt.Printf("-%v\n", missingFlag)
		}
		return false
	}

	return true
}

// Rebuilds the command line arguments, excluding the daemon flag
func RebuildArgs(defined map[string]string, daemonFlag string) []string {
	rebuiltLength := len(defined)
	_, ok := defined[daemonFlag]
	if ok {
		rebuiltLength--
	} else {
		fmt.Printf("Warning daemon flag (-%v) not found\n", daemonFlag)
	}
	rebuiltArgs := make([]string, rebuiltLength)
	argIndex := 0
	for argKey, argValue := range defined {
		if argKey != daemonFlag {
			rebuiltArgs[argIndex] = fmt.Sprintf("-%v=%v", argKey, argValue)
			argIndex++
		}
	}
	return rebuiltArgs
}

// Set up the various loggers
func PrepareLogging(logFile *os.File, debug bool) (main, heartbeat, launch, script *log.Logger) {
	mainLogOut := io.Discard
	hbLogOut := io.Discard
	launchLogOut := io.Discard
	if debug {
		mainLogOut = logFile
		hbLogOut = logFile
		launchLogOut = logFile
	}
	mainLogger := log.New(mainLogOut, "MAIN ", log.Lmicroseconds|log.Lshortfile)
	hbLogger := log.New(hbLogOut, "HEARBEAT ", log.Lmicroseconds|log.Lshortfile)
	launchLogger := log.New(launchLogOut, "LAUNCHER ", log.Lmicroseconds|log.Lshortfile)
	scriptLogger := log.New(logFile, "", log.Lmicroseconds|log.Lshortfile)

	return mainLogger, hbLogger, launchLogger, scriptLogger
}

// Catch termination signals to allow for a graceful exit (i.e. no zombies)
// Only for this process, does not catch any signals to the launched script.
func SignalCatcher(sigChan chan os.Signal, logger *log.Logger) {
	for sig := range sigChan {
		logger.Printf("(sig catcher) caught: %v\n", sig)
	}
}

// Touches log file while launched script is still active
func Heartbeat(wg *sync.WaitGroup, exitChan chan bool,
	controlDir string, resultPath string, logPath string, logger *log.Logger) {

	defer wg.Done()
	defer close(exitChan)

	_, err := os.Stat(controlDir)
	if os.IsNotExist(err) {
		logger.Printf("%v\n", err.Error())
		return
	}
	_, err = os.Stat(resultPath)
	if !os.IsNotExist(err) {
		logger.Printf("Result file already exists, stopping heartbeat.\n%v\n", resultPath)
		return
	}

	for {
		select {
		case <-exitChan:
			logger.Println("received script finished, exiting")
			return
		default:
			// heartbeat
			logger.Println("touch log")
			err = os.Chtimes(logPath, time.Now(), time.Now())
			CheckIfErr(logger, err)

			time.Sleep(time.Second * 3)
		}
	}
}

// Write launched script's exit code to a file
func RecordExit(exitCode int, resultPath string, logger *log.Logger) {
	resultFile, err := os.Create(resultPath)
	if CheckIfErr(logger, err) {
		return
	}
	defer resultFile.Close()
	_, err = resultFile.WriteString(strconv.Itoa(exitCode))
	CheckIfErr(logger, err)
	err = resultFile.Close()
	CheckIfErr(logger, err)
}
