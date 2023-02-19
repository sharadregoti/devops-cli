package utils

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"github.com/kr/pty"
	"github.com/sharadregoti/devops/utils/logger"
	"golang.org/x/term"
)

func GetTableTitle(title string, total int) string {
	return fmt.Sprintf("%s (%d)", title, total)
}

func ParseTableTitle(title string) (string, string) {
	arr := strings.Split(title, " ")
	if len(arr) == 0 {
		return "", ""
	}
	if len(arr) == 1 {
		return arr[0], ""
	}
	return strings.ToLower(arr[0]), strings.TrimPrefix(strings.TrimSuffix(arr[1], ")"), "(")
}

func ExecuteCMDGetOutput(cmdStr string) (string, error) {
	// create a temporary file
	f, err := ioutil.TempFile("", "script")
	if err != nil {
		// handle error
		return "", logger.LogError("failed to create temp file: Error ", err)
	}
	defer os.Remove(f.Name()) // delete the file when we're done

	// write the script to the file
	_, err = f.WriteString(cmdStr)
	if err != nil {
		// handle error
		f.Close()
		return "", logger.LogError("failed to write data in temp file: Error ", err)
	}
	f.Close()

	// make the file executable
	err = os.Chmod(f.Name(), 0777)
	if err != nil {
		// handle error
		return "", logger.LogError("failed to make the temp file executable: Error ", err)
	}

	// Set the command to execute
	// Execute the command
	output, err := exec.Command(f.Name()).Output()
	if err != nil {
		return "", logger.LogError("Error while executing command\n %v\n: %s", cmdStr, err)
	}

	return string(output), nil
}

func ExecuteCMDLong(cmdStr string) error {
	// Set the command to execute
	arr := strings.Split(cmdStr, " ")

	// Execute the command
	err := exec.Command(arr[0], arr[1:]...).Start()
	if err != nil {
		return logger.LogError("Error while executing command: %s", err)
	}

	return nil
}

func ExecuteCMD(cmdStr string) error {
	// Create arbitrary command.

	// arr := strings.Split(cmdStr, " ")
	// create a temporary file
	f, err := ioutil.TempFile("", "script")
	if err != nil {
		// handle error
		return logger.LogError("failed to create temp file: Error ", err)
	}
	defer os.Remove(f.Name()) // delete the file when we're done

	// write the script to the file
	_, err = f.WriteString(cmdStr)
	if err != nil {
		// handle error
		f.Close()
		return logger.LogError("failed to write data in temp file: Error ", err)
	}
	f.Close()

	// make the file executable
	err = os.Chmod(f.Name(), 0777)
	if err != nil {
		// handle error
		return logger.LogError("failed to make the temp file executable: Error ", err)
	}

	c := exec.Command(f.Name())

	// Start the command with a pty.
	ptmx, err := pty.Start(c)
	if err != nil {
		return err
	}

	// Make sure to close the pty at the end.
	defer func() {
		_ = ptmx.Close()
		logger.LogDebug("Closing pty")
	}() // Best effort.

	// Handle pty size.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)
	go func() {
		for range ch {
			if err := pty.InheritSize(os.Stdin, ptmx); err != nil {
				log.Printf("error resizing pty: %s", err)
			}
		}
		logger.LogDebug("Closing notify channel")
	}()
	ch <- syscall.SIGWINCH // Initial resize.
	defer func() {
		signal.Stop(ch)
		close(ch)
		logger.LogDebug("Closing signal channel")
	}() // Cleanup signals when done.

	// Set stdin in raw mode.
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = term.Restore(int(os.Stdin.Fd()), oldState)
		logger.LogDebug("Closing make raw chan")
	}() // Best effort.

	// Copy stdin to the pty and the pty to stdout.
	// NOTE: The goroutine will keep reading until the next keystroke before returning.
	go func() {
		_, _ = io.Copy(ptmx, os.Stdin)
		logger.LogDebug("Closing stdin:")
	}()

	// TODO: Errors made by commands are not being handled
	_, _ = io.Copy(os.Stdout, ptmx)
	logger.LogDebug("Closing stdout:")
	logger.LogDebug("Exit code is: %d", c.ProcessState.ExitCode())
	return nil
}
