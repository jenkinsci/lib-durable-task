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

package common_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"jenkinsci.org/plugins/durabletask/common"
)

var argMap = map[string]string{
	"a": "0",
	"b": "1",
	"c": "2",
	"d": "3",
	"e": "4",
	"f": "5",
	"g": "6",
	"h": "7",
	"i": "8",
	"j": "9",
}

func TestValidateFlags(t *testing.T) {
	required := []string{"a", "j", "f"}
	result := common.ValidateFlags(argMap, required)
	if result != true {
		t.Error("expected true")
	}

	requireNone := []string{}
	result = common.ValidateFlags(argMap, requireNone)
	if result != true {
		t.Error("expected true")
	}

	requiredDupes := []string{"a", "j", "j", "b", "a"}
	result = common.ValidateFlags(argMap, requiredDupes)
	if result != true {
		t.Error("expected true")
	}
}

func ExampleValidateFlags_missing() {
	requiresMore := []string{"a", "j", "z", "p", "m"}
	result := common.ValidateFlags(argMap, requiresMore)
	if result != false {
		log.Fatal("expected false")
	}

	// Output:
	// The following required flags are missing:
	// -z
	// -p
	// -m
}

func TestRebuildArgs(t *testing.T) {
	removedFlag := "d"
	removedArg := "-" + removedFlag + "=" + argMap[removedFlag]
	rebuiltArgs := common.RebuildArgs(argMap, removedFlag)
	for i := range rebuiltArgs {
		if rebuiltArgs[i] == removedArg {
			t.Errorf("flag %v was not removed correctly: %v", removedFlag, rebuiltArgs)
		}
	}

	r, w, err := os.Pipe()
	if err != nil {
		t.Error(err)
	}
	origStdout := os.Stdout
	os.Stdout = w

	removedFlag = "z"
	rebuiltArgs = common.RebuildArgs(argMap, removedFlag)
	if len(rebuiltArgs) != len(argMap) {
		t.Errorf("non-existent flag %v incorrectly modified args: %v", removedFlag, rebuiltArgs)
	}

	buf := make([]byte, 1024)
	n, err := r.Read(buf)
	if err != nil {
		t.Error(err)
	}
	os.Stdout = origStdout

	warningMsg := string(buf[:n])
	if warningMsg != "Warning daemon flag (-z) not found\n" {
		t.Errorf("Warning to stdout does not match: %v", warningMsg)
	}
}

func TestHeartbeatAndSignalFinished(t *testing.T) {
	currentDir, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}

	logPath := filepath.FromSlash(currentDir + "/logfile")
	logFile, err := os.Create(filepath.FromSlash(currentDir + "/logfile"))
	if err != nil {
		t.Error(err)
	}
	defer func() {
		logFile.Close()
		os.Remove(logFile.Name())
	}()

	resultPath := filepath.FromSlash(currentDir + "/resultfile")
	// make sure no old result files remain
	os.Remove(resultPath)

	var buf bytes.Buffer
	hbLogger := log.New(&buf, "", 0)

	exitChan := make(chan bool)
	var wg sync.WaitGroup
	wg.Add(1)
	fmt.Println("launching heartbeat")
	go common.Heartbeat(&wg, exitChan, currentDir, resultPath, logPath, hbLogger)

	fmt.Println("sleeping for 6 seconds")
	time.Sleep(6 * time.Second)
	fmt.Println("signal finished to heartbeat")
	common.SignalFinished(exitChan)
	fmt.Print("waiting...")
	wg.Wait()
	fmt.Println("done.")

	output := buf.String()
	if !strings.Contains(output, "touch log\ntouch log\n") {
		t.Errorf("Output did not contain at least 2 'touch log' statements:\n%v", output)
	}
	if !strings.Contains(output, "received script finished, exiting") {
		t.Errorf("Ouput did not contain signal finished statement:\n%v", output)
	}
}

func TestHeartbeatResultFound(t *testing.T) {
	currentDir, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}

	logPath := filepath.FromSlash(currentDir + "/logfile")
	logFile, err := os.Create(filepath.FromSlash(currentDir + "/logfile"))
	if err != nil {
		t.Error(err)
	}
	defer func() {
		logFile.Close()
		os.Remove(logFile.Name())
	}()

	resultPath := filepath.FromSlash(currentDir + "/resultfile")
	resultFile, err := os.Create(filepath.FromSlash(currentDir + "/resultfile"))
	if err != nil {
		t.Error(err)
	}
	defer func() {
		resultFile.Close()
		os.Remove(resultFile.Name())
	}()

	var buf bytes.Buffer
	hbLogger := log.New(&buf, "", 0)

	exitChan := make(chan bool)
	var wg sync.WaitGroup
	wg.Add(1)
	go common.Heartbeat(&wg, exitChan, currentDir, resultPath, logPath, hbLogger)
	wg.Wait()

	output := buf.String()
	expected := fmt.Sprintf("Result file already exists, stopping heartbeat.\n%v\n", resultPath)
	if output != expected {
		t.Errorf("Mismatched output.\nWanted:\n%v\nActual:\n%v", expected, output)
	}
}

func TestRecordExit(t *testing.T) {
	currentDir, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}
	resultPath := filepath.FromSlash(currentDir + "/resultfile")
	// make sure no old result files remain
	os.Remove(resultPath)

	common.RecordExit(-12345678, resultPath, log.New(os.Stdout, "", 0))
	data, err := ioutil.ReadFile(resultPath)
	if err != nil {
		t.Error(err)
	}
	if string(data) != "-12345678" {
		t.Errorf("result data does not match. Expected -12345678, got %v", string(data))
	}

	// Leave old result file present
	common.RecordExit(99999, resultPath, log.New(os.Stdout, "", 0))
	data, err = ioutil.ReadFile(resultPath)
	if err != nil {
		t.Error(err)
	}
	if string(data) != "99999" {
		t.Errorf("result data does not match. Expected 99999, got %v", string(data))
	}
	os.Remove(resultPath)
}
