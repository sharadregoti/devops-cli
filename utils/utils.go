package utils

import (
	"fmt"
	"io"
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
	return strings.ToLower(arr[0]), strings.TrimPrefix(strings.TrimSuffix(arr[1], ")"), "(")
}

func ExecuteCMDGetOutput(cmdStr string) (string, error) {
	// Set the command to execute
	arr := strings.Split(cmdStr, " ")

	// Execute the command
	output, err := exec.Command(arr[0], arr[1:]...).Output()
	if err != nil {
		return "", logger.LogError("Error while executing command %v: %s", arr, err)
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

	arr := strings.Split(cmdStr, " ")

	c := exec.Command(arr[0], arr[1:]...)

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
	return nil
}
