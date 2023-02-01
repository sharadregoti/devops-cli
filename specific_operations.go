package core

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"github.com/google/uuid"
	"github.com/kr/pty"
	"github.com/sharadregoti/devops/model"
	"github.com/sharadregoti/devops/shared"
	"github.com/sharadregoti/devops/utils"
	"github.com/sharadregoti/devops/utils/logger"
	"golang.org/x/term"
)

func (c *CurrentPluginContext) ReadFromSTDOUT(ch chan string) error {
	newReader := bufio.NewReader(c.pc.GetStdoutReader())
	if newReader == nil {
		return logger.LogError("Failed to create a reader for stream")
	}

	go func() {
		for {
			logger.LogDebug("Started listener for ReadFromSTDOUT")
			// select {
			// case <-streamCloserChan:
			// 	logger.LogDebug("Streamer go routine closed")
			// 	return
			// default:
			// newReader.ReadString()
			// data, _, err := newReader.ReadLine()
			// newReader.read
			data, err := newReader.ReadByte()
			if err == io.EOF {
				logger.LogError("EOF received while streaming: %v", err)
				break
			} else if err != nil {
				logger.LogError("failed to stream data: %v", err)
				break
			}
			logger.LogDebug("Reading some some data: (%s)", string(data))
			ch <- string(data)
			logger.LogDebug("Sending read finished")
			// }
			// logger.LogDebug("Closed listener for ReadFromSTDOUT")
		}
	}()

	return nil
}

func (c *CurrentPluginContext) WriteIntoSTDIN(ch chan string) error {
	// newWriter := bufio.NewWriter(c.pc.GetStdoutWriter())
	// if newWriter == nil {
	// 	return logger.LogError("Failed to create a reader for stream")
	// }

	// newWriter := bufio.NewWriter(os.Stdout)
	// if newWriter == nil {
	// 	return logger.LogError("Failed to create a reader for stream")
	// }

	go func() {
		logger.LogDebug("Started listener for WriteIntoSTDIN")
		for d := range ch {
			logger.LogDebug("Writing some data: (%s)", string(d))
			for _, b := range d {
				_, err := fmt.Fprint(os.Stdin, byte(b))
				// err := newWriter.WriteByte(byte(b))
				if err != nil {
					logger.LogError("failed to write into stdout: %v", err.Error())
				}
				// err = newWriter.Flush()
				// if err != nil {
				// 	logger.LogError("failed to flush into stdout: %v", err.Error())
				// }
			}
			logger.LogDebug("Writing Finished")
		}
		logger.LogDebug("Closing listener for WriteIntoSTDIN")
	}()

	return nil
}

func (c *CurrentPluginContext) PerformSavedAction(id string, rw io.ReadWriter) error {
	a, ok := c.actionsToExecute[id]
	if !ok {
		return fmt.Errorf("id (%s) does not exists in saved action map", id)
	}

	fnArgs := shared.SpecificActionArgs{
		ActionName:   a.e.Type,
		ResourceName: a.e.ResourceName,
		ResourceType: a.e.ResourceType,
		IsolatorName: a.e.IsolatorName,
		Args:         nil,
	}

	logger.LogDebug("Performing specific action: %v", fnArgs)

	res, err := c.plugin.PerformSpecificAction(fnArgs)
	if err != nil {
		return err
	}

	cmd := res.Result.(string)
	logger.LogDebug("Performing specific action got result: %v", cmd)

	if err := cmdExec(cmd, rw); err != nil {
		return err
	}

	return nil
}

func cmdExec(cmdStr string, wr io.ReadWriter) error {
	// Create arbitrary command.
	arr := strings.Split(cmdStr, " ")

	c := exec.Command(arr[0], arr[1:]...)

	// Start the command with a pty.
	ptmx, err := pty.Start(c)
	if err != nil {
		return err
	}

	// Make sure to close the pty at the end.
	defer func() { _ = ptmx.Close() }() // Best effort.

	// Handle pty size.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)
	go func() {
		for range ch {
			if err := pty.InheritSize(os.Stdin, ptmx); err != nil {
				log.Printf("error resizing pty: %s", err)
			}
		}
	}()
	ch <- syscall.SIGWINCH                        // Initial resize.
	defer func() { signal.Stop(ch); close(ch) }() // Cleanup signals when done.

	// Set stdin in raw mode.
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer func() { _ = term.Restore(int(os.Stdin.Fd()), oldState) }() // Best effort.

	// Copy stdin to the pty and the pty to stdout.
	// NOTE: The goroutine will keep reading until the next keystroke before returning.
	go func() { _, _ = io.Copy(ptmx, wr) }()

	_, _ = io.Copy(wr, ptmx)

	return nil
}

func (c *CurrentPluginContext) saveAction(e model.Event) (*model.EventResponse, error) {
	id := uuid.NewString()
	logger.LogDebug("Saving action (%s) againts id (%s)", e.Type, id)
	c.actionsToExecute[id] = &actionsToExecute{e: e}
	return &model.EventResponse{ID: id}, nil
}

func (c *CurrentPluginContext) ExecuteSpecificActionTemplate(a shared.Action, e model.Event) (string, error) {
	res, err := c.Read(e)
	if err != nil {
		return "", logger.LogError("failed to read resource: %v", err)
	}

	params := map[string]interface{}{
		"resourceName": e.ResourceName,
		"resourceType": e.ResourceType,
		"isolatorName": e.IsolatorName,
		"resource":     res,
		"args":         e.Args,
	}

	templateExecutedArgs := map[string]interface{}{}
	for key, value := range e.Args {
		strValue, ok := value.(string)
		if ok && strValue != "" {
			logger.LogDebug("Executing template args having key (%s)", key)
			res, err := utils.ExecuteTemplate(strValue, params)
			if err != nil {
				return "", logger.LogError("failed to execute template: %v", err)
			}
			templateExecutedArgs[key] = res
			continue
		}
		templateExecutedArgs[key] = value
	}

	tempRes, err := utils.ExecuteTemplate(a.Execution.Cmd, params)
	if err != nil {
		return "", logger.LogError("failed to execute template: %v", err)
	}

	return tempRes, nil
}

func (c *CurrentPluginContext) SpecificAction(a shared.Action, e model.Event) (*model.EventResponse, error) {
	fnArgs := shared.SpecificActionArgs{
		ActionName:   e.Type,
		ResourceName: e.ResourceName,
		ResourceType: e.ResourceType,
		IsolatorName: e.IsolatorName,
		Args:         e.Args,
	}

	var result interface{}
	var err error
	if a.Execution.Cmd != "" {
		result, err = c.ExecuteSpecificActionTemplate(a, e)
		if err != nil {
			return nil, err
		}
	} else {
		// Execute actual action
		res, err := c.plugin.PerformSpecificAction(fnArgs)
		if err != nil {
			return nil, err
		}
		result = res.Result
	}

	switch a.OutputType {

	case string(model.OutputTypeString):
		res, err := utils.ExecuteCMDGetOutput(result.(string))
		if err != nil {
			return nil, err
		}
		return &model.EventResponse{Result: res}, nil

	case string(model.OutputTypeNothing):
		if a.Execution.IsLongRunning {
			err := utils.ExecuteCMDLong(result.(string))
			if err != nil {
				return nil, err
			}
		} else {
			_, err := utils.ExecuteCMDGetOutput(result.(string))
			if err != nil {
				return nil, err
			}
		}
		return &model.EventResponse{}, nil

	case string(model.OutputTypeBidrectional), string(model.OutputTypeStream):
		// res, err := c.plugin.PerformSpecificAction(fnArgs)
		// if err != nil {
		// 	return nil, err
		// }
		// TODO: remove action after some time
		saRes, err := c.saveAction(e)
		if err != nil {
			return nil, err
		}

		return &model.EventResponse{Result: result, ID: saRes.ID}, nil

	default:
		return nil, fmt.Errorf("invalid output type (%v) provided for executing specfic action", a.OutputType)
	}
}
