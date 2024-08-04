package client

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/ByteStorage/FlyDB/db/grpc/service"
)

var (
	stdoutPipe     io.ReadCloser
	stdoutScanner  *bufio.Scanner
	stderr         = new(bytes.Buffer)
	stdin          io.WriteCloser
	stdoutStringCh = make(chan string)
	currentIndex   = 0
)

func checkStdoutTextIsEmptyOrNot() {
	for {
		select {
		case stdoutText := <-stdoutStringCh:
			if stdoutText != "" {
				panic(fmt.Sprintf("stdout text is not empty: %s", stdoutText))
			}
		case <-time.After(1 * time.Second):
			return
		}
	}
}

func closeFlyDBClient(cmd *exec.Cmd) func() {
	return func() {
		var err error

		checkStdoutTextIsEmptyOrNot()

		close(stdoutStringCh)

		if err = stdin.Close(); err != nil {
			fmt.Println("close std in err: ", stderr)
			panic(err)
		}

		if err = cmd.Wait(); err != nil {
			fmt.Println("cmd wait err: ", stderr)
			panic(err)
		}
	}
}

func setupStdinPipe(cmd *exec.Cmd) {
	var (
		err error
	)

	if stdin, err = cmd.StdinPipe(); err != nil {
		panic(err)
	}
}

func setupStdout(cmd *exec.Cmd) {
	var (
		err error
	)

	if stdoutPipe, err = cmd.StdoutPipe(); err != nil {
		panic(err)
	}

	stdoutScanner = bufio.NewScanner(stdoutPipe)
}

func cmdStart(cmd *exec.Cmd) {
	var (
		err error
	)

	if err = cmd.Start(); err != nil {
		panic(err)
	}
}

func checkFlyDBOutput(currentText string) bool {
	ignoreMessages := []string{
		"        ___             ___        ___             ___             ___  ",
		"       /\\  \\           /\\__\\      /\\__\\           /\\  \\           /\\  \\  ",
		"      /::\\  \\         /:/  /     /:/  / ___      /::\\  \\         /::\\  \\  ",
		"     /:/\\ \\  \\       /:/  /     /:/  / /\\__\\    /:/\\ \\  \\       /:/\\ \\  \\  ",
		"    /::\\~\\ \\  \\     /:/  /      \\:\\  \\/ /  /   /:/ /\\ \\  \\     /::\\ \\ \\  \\  ",
		"   /:/\\:\\ \\ \\__\\   /:/  /      __\\:\\~/ /  /   /:/_/  \\ \\  \\   /:/\\:\\ \\ \\  \\  ",
		"   \\/__\\:\\ \\/__/   \\:\\  \\     /\\  \\:::/  /    \\:\\ \\  | |  |   \\:\\ \\:\\/ |  |  ",
		"        \\:\\__\\      \\:\\  \\    \\:\\~~/:/  /      \\:\\~\\/ /  /     \\:\\ \\::/  /  ",
		"         \\/__/       \\:\\  \\    \\:\\/:/  /        \\:\\/ /  /       \\:\\/:/  /  ",
		"                      \\ \\__\\    \\::/  /          \\::/  /         \\::/  /  ",
		"                       \\/__/     \\/__/            \\/__/           \\/__/                                                                        ",
	}

	if currentIndex < len(ignoreMessages) && currentText == ignoreMessages[currentIndex] {
		currentIndex++
		return true
	}

	return false
}

func startFlyDBClient() func() {
	cmd := exec.Command("go", "run", "./cli/flydb-client.go", Addr)

	cmd.Stderr = stderr

	setupStdinPipe(cmd)

	setupStdout(cmd)

	go func() {

		for stdoutScanner.Scan() {
			currentText := stdoutScanner.Text()

			if currentText == "" {
				continue
			}
			if checkFlyDBOutput(currentText) {
				continue
			}

			stdoutStringCh <- stdoutScanner.Text()
		}
	}()

	cmdStart(cmd)

	return closeFlyDBClient(cmd)
}

type cmdMeta struct {
	cmd          string
	expectResult []string
	testcaseName string
}

func (c *cmdMeta) lofFileAndLine(t *testing.T, funcName string) {
	found := false
	for skip := 0; !found; skip++ {
		pc, file, line, ok := runtime.Caller(skip)
		if !ok {
			break
		}
		details := runtime.FuncForPC(pc)
		if details != nil && strings.Contains(details.Name(), funcName) {
			t.Logf("TestCase : %s in %s:%d", funcName, file, line)
			found = true
		}
	}
	if !found {
		t.Logf("File: %s, Line: %d", "unknown", 0)
	}
}

func (c *cmdMeta) cmdTest(t *testing.T) {
	var (
		err          error
		stdoutString string
	)

	if _, err = stdin.Write([]byte(c.cmd + "\n")); err != nil {
		t.Fatal(err)
	}

	for _, expectResult := range c.expectResult {
		select {
		case stdoutString = <-stdoutStringCh:
		case <-time.After(1 * time.Second):
			c.lofFileAndLine(t, c.testcaseName)
			t.Fatal("not get stdout")
		}

		if strings.Contains(stdoutString, expectResult) == false {
			c.lofFileAndLine(t, c.testcaseName)
			t.Fatalf("expect: %s\n got: %s\n", expectResult, stdoutString)
		}
	}

}

func TestMain(m *testing.M) {
	var (
		ctx                       = context.Background()
		err                       error
		serviceContainer          = &service.Container{}
		closeServiceContainerFunc func()
		closeFlyDBClientFunc      func()
	)

	if closeServiceContainerFunc, err = serviceContainer.StartServiceContainer(ctx, flyDBServerPort); err != nil {
		panic(err)
	}

	Addr = serviceContainer.URI

	closeFlyDBClientFunc = startFlyDBClient()

	code := m.Run()

	closeFlyDBClientFunc()
	closeServiceContainerFunc()

	os.Exit(code)
}
