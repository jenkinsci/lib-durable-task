package main

import (
	"bytes"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"sync"
	"syscall"
	"testing"
	"time"

	"golang.org/x/sys/unix"
	"jenkinsci.org/plugins/durabletask/common"
)

// Test SignalCatcher here instead of common test since it's OS dependent
func ExampleSignalCatcher() {
	logger := log.New(os.Stdout, "", 0)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, unix.SIGINT, unix.SIGTERM, unix.SIGHUP)

	go common.SignalCatcher(sigChan, logger)
	err := syscall.Kill(os.Getpid(), unix.SIGINT)
	if err != nil {
		log.Fatal(err.Error())
	}
	time.Sleep(100 * time.Millisecond)
	err = syscall.Kill(os.Getpid(), unix.SIGTERM)
	if err != nil {
		log.Fatal(err.Error())
	}
	time.Sleep(100 * time.Millisecond)
	err = syscall.Kill(os.Getpid(), unix.SIGHUP)
	if err != nil {
		log.Fatal(err.Error())
	}
	time.Sleep(100 * time.Millisecond)
	signal.Stop(sigChan)
	close(sigChan)

	// Output:
	// (sig catcher) caught: interrupt
	// (sig catcher) caught: terminated
	// (sig catcher) caught: hangup
}

func TestLauncher(t *testing.T) {
	cookieName := "aCookie"
	cookieVal := "withRaisins"
	interpreter := "sh"
	currentDir, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}
	resultPath := filepath.FromSlash(currentDir + "/resultfile")
	outputPath := filepath.FromSlash(currentDir + "/outputfile")
	// clear old artifacts
	os.Remove(resultPath)
	os.Remove(outputPath)
	scriptPath := filepath.FromSlash(currentDir + "/test-script.sh")

	var launchBuffer bytes.Buffer
	launchLogger := log.New(&launchBuffer, "", 0)
	var scriptBuffer bytes.Buffer
	scriptLogger := log.New(&scriptBuffer, "", 0)

	exitChan := make(chan bool)
	var wg sync.WaitGroup
	wg.Add(1)
	go launcher(&wg, exitChan, cookieName, cookieVal,
		interpreter, scriptPath, resultPath, outputPath,
		launchLogger, scriptLogger)
	// Must access the channel else it blocks launcher
	channelResult := <-exitChan
	if channelResult != true {
		t.Errorf("Channel signaled incorrectly: %v", channelResult)
	}
	wg.Wait()

	launchOutput := launchBuffer.String()
	launchLoggerExp := regexp.MustCompile(`^args 0: sh\nargs 1: -xe\nargs 2: .*\/cmd\/bash\/test-script.sh\nlaunched \d*\nscript exit code: 0`)
	if !launchLoggerExp.MatchString(launchOutput) {
		t.Errorf("launch output incorrect:\n%v", launchOutput)
	}
	scriptOutput := scriptBuffer.String()
	if scriptOutput != "+ echo hello\n" {
		t.Errorf("script output incorrect:\n%v", scriptOutput)
	}
	os.Remove(resultPath)
	os.Remove(outputPath)
}
