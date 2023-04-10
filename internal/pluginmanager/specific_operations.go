package pluginmanager

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
	"github.com/kr/pty"
	"github.com/sharadregoti/devops-plugin-sdk/proto"
	"github.com/sharadregoti/devops/model"
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

	fnArgs := &proto.SpecificActionArgs{
		ActionName:   a.e.Type,
		ResourceName: a.e.ResourceName,
		ResourceType: a.e.ResourceType,
		IsolatorName: a.e.IsolatorName,
		Args:         nil,
	}

	logger.LogDebug("Performing specific action: %v", fnArgs)

	// res, err := c.plugin.PerformSpecificAction(fnArgs)
	// if err != nil {
	// 	return err
	// }

	// cmd := res.Result.AsInterface().(string)
	logger.LogDebug("Performing specific action got result: %v", a.cmd)

	if err := cmdExec(a.cmd, rw); err != nil {
		return err
	}

	return nil
}

func cmdExec(cmdStr string, wr io.ReadWriter) error {
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

func (c *CurrentPluginContext) saveAction(cmd string, e model.Event) (*model.EventResponse, error) {
	id := uuid.NewString()
	logger.LogDebug("Saving action (%s) againts id (%s)", e.Type, id)
	c.actionsToExecute[id] = &actionsToExecute{e: e, cmd: cmd}
	return &model.EventResponse{ID: id}, nil
}
func (c *CurrentPluginContext) ExecuteSpecificActionTemplate(a *proto.Action, e model.Event) (string, error) {

	params, err := c.getTemplateParams(e)
	if err != nil {
		return "", err
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

func (c *CurrentPluginContext) getTemplateParams(e model.Event) (map[string]interface{}, error) {
	res, err := c.Read(e)
	if err != nil {
		return nil, logger.LogError("failed to read resource: %v", err)
	}

	params := map[string]interface{}{
		"resourceName": e.ResourceName,
		"resourceType": e.ResourceType,
		"isolatorName": e.IsolatorName,
		"authPath":     c.authInfo.Path,
		"authName":     c.authInfo.Name,
		"authId":       c.authInfo.IdentifyingName,
		"resource":     res,
		"args":         e.Args,
	}

	return params, nil
}

func (c *CurrentPluginContext) ExecuteSpecificActionTemplateArgs(e model.Event) (map[string]interface{}, error) {

	params, err := c.getTemplateParams(e)
	if err != nil {
		return nil, err
	}

	templateExecutedArgs := map[string]interface{}{}
	for key, value := range e.Args {
		strValue, ok := value.(string)
		if ok && strValue != "" {
			logger.LogDebug("Executing template args having key (%s)", key)
			res, err := utils.ExecuteTemplate(strValue, params)
			if err != nil {
				return nil, logger.LogError("failed to execute template: %v", err)
			}
			templateExecutedArgs[key] = res
			continue
		}
		templateExecutedArgs[key] = value
	}

	return templateExecutedArgs, nil
}

func (c *CurrentPluginContext) SpecificAction(a *proto.Action, e model.Event) (*model.EventResponse, error) {
	fnArgs := &proto.SpecificActionArgs{
		ActionName:   e.Type,
		ResourceName: e.ResourceName,
		ResourceType: e.ResourceType,
		IsolatorName: e.IsolatorName,
		Args:         utils.GetMap(e.Args),
	}

	var result interface{}
	var err error
	if a.Execution.Cmd != "" {
		logger.LogDebug("Execting template...")
		result, err = c.ExecuteSpecificActionTemplate(a, e)
		if err != nil {
			return nil, err
		}
		logger.LogDebug("Command template result is (%s)", result)
	} else {
		logger.LogDebug("Execting actual action...")
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
			// err := utils.ExecuteCMDLong(result.(string))
			err := c.NewLongRunning(result.(string), &e)
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
		saRes, err := c.saveAction(result.(string), e)
		if err != nil {
			return nil, err
		}

		return &model.EventResponse{Result: result, ID: saRes.ID}, nil

	default:
		return nil, fmt.Errorf("invalid output type (%v) provided for executing specfic action", a.OutputType)
	}
}

func (c *CurrentPluginContext) NewLongRunning(cmdStr string, e *model.Event) error {
	id := uuid.NewString()
	lri := &model.LongRunningInfo{
		ID:      id,
		Name:    e.Type,
		Status:  "running",
		Message: "NA",
	}
	lri.SetE(e)

	// create a temporary file
	f, err := ioutil.TempFile("", "script")
	if err != nil {
		// handle error
		return logger.LogError("failed to create temp file: Error ", err)
	}

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

	// Execute the command
	cmd := exec.Command(f.Name())
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true, Pgid: 0}
	lri.SetCMD(cmd)
	if err := cmd.Start(); err != nil {
		return logger.LogError("Error while running command: %s", err)
	}

	go func() {
		defer logger.LogDebug("Exiting long running command routine with id (%s)", id)
		defer os.Remove(f.Name()) // delete the file when we're done

		err := cmd.Wait()
		if err != nil {
			lri.Status = "failed"
			lri.Message = err.Error()
			logger.LogError("Error while waiting for command to finish: %s", err)
		}
	}()

	c.longRunning[id] = lri
	logger.LogDebug("Added long running command with id (%s), process id (%s)", id, cmd.Process.Pid)
	return nil
}

// func ExecuteCMDLong(cmdStr string) error {
// 	// Set the command to execute
// 	arr := strings.Split(cmdStr, " ")

// 	// Execute the command
// 	c := exec.Command(arr[0], arr[1:]...)
// 	err := c.Run()
// 	if err != nil {
// 		return logger.LogError("Error while running command: %s", err)
// 	}

// 	go func() {
// 		err := c.Wait()
// 		if err != nil {
// 			logger.LogError("Error while waiting for command to finish: %s", err)
// 		}

// 	}()

// 	return nil
// }
