package main

import (
	"log"
	"syscall"
)

func sendSignal(pid int) {
	d, e := syscall.LoadDLL("kernel32.dll")
	if e != nil {
		log.Fatalf("LoadDLL: %v\n", e)
	}
	p, e := d.FindProc("GenerateConsoleCtrlEvent")
	if e != nil {
		log.Fatalf("FindProc: %v\n", e)
	}
	r, _, e := p.Call(syscall.CTRL_C_EVENT, uintptr(pid))
	if r == 0 {
		log.Fatalf("GenerateConsoleCtrlEvent: %v\n", e)
	}
}

// windows interrupt signals cannot be sent from a program.
// see: github.com/golang/go/issues/6720
// IMPORTANT: This test is commented out because it will fail for the following reasons:
//   1) will ALWAYS fail when executed as RUN command in dockerfile (even on > ltsc2019)
//   2) code only works on windows versions > ltsc2019 (https://github.com/docker/for-win/issues/3173)
/*
func ExampleSignalCatcher() {
	logger := log.New(os.Stdout, "", 0)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, windows.SIGINT, windows.SIGTERM)

	go common.SignalCatcher(sigChan, logger)

	const source = `
package main
import (
	"log"
	"os"
	"os/signal"
	"time"

	"golang.org/x/sys/windows"
	"jenkinsci.org/plugins/durabletask/common"
)
func main() {
	logger := log.New(os.Stdout, "", 0)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, windows.SIGINT, windows.SIGTERM)

	go common.SignalCatcher(sigChan, logger)

	time.Sleep(3 * time.Second)
	signal.Stop(sigChan)
	close(sigChan)
}
`
	tmp, err := os.MkdirTemp("", "sigcatcher")
	if err != nil {
		log.Fatal(err)
	}
	name := filepath.Join(tmp, "catcher")
	src := name + ".go"
	exe := name + ".exe"

	// write ctrlbreak.go
	f, err := os.Create(src)
	if err != nil {
		log.Fatalf("Failed to create %v: %v", src, err)
	}
	defer f.Close()
	f.Write([]byte(source))

	// compile it
	defer os.Remove(exe)
	o, err := exec.Command("go", "build", "-o", exe, src).CombinedOutput()
	if err != nil {
		log.Fatalf("Failed to compile: %v\n%v", err, string(o))
	}

	// run it
	cmd := exec.Command(exe)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
	}
	err = cmd.Start()
	if err != nil {
		log.Fatalf("Start failed: %v", err)
	}
	go func() {
		time.Sleep(500 * time.Millisecond)
		sendSignal(os.Getpid())
	}()
	err = cmd.Wait()
	if err != nil {
		log.Fatalf("Program exited with error: %v\n", err)
	}

	// Output:
	// (sig catcher) caught: interrupt
}
*/
