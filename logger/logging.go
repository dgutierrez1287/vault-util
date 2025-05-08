package logger

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/hashicorp/go-hclog"
)

var Logger hclog.Logger
var LogLevel string
var machineOutput bool

/*
initialize logging, this will set level, colorization
and set machine output setting for logger
*/
func InitLogging(debug bool, colorize bool, machineOnlyOutput bool) {
  // set up logger
  machineOutput = machineOnlyOutput
  var colorOpt hclog.ColorOption

  // debug output setup
  if debug {
    LogLevel = "DEBUG"
    fmt.Println("Debugging enabled")
  } else {
    LogLevel = "INFO"
  }

  // colorization setup 
  if colorize {
    colorOpt = hclog.ColorOption(hclog.AutoColor)
  } else {
    colorOpt = hclog.ColorOption(hclog.ColorOff)
  }

  if !machineOnlyOutput {
    fmt.Printf("loglevel %s \n", LogLevel)
  }

  // create global logger
  Logger = hclog.New(&hclog.LoggerOptions{
    Name: "vault-util",
    Level: hclog.LevelFromString(LogLevel),
    Color: colorOpt,
  })
}

/*
This will wrap info logging to handle if it should be 
written to console based on machine output setting
*/
func LogInfo(message string, args ...interface{}) {
  if !machineOutput {
    Logger.Info(message, args...)
  }
}

/*
This will wrap debug logging to handle any special 
caes
*/
func LogDebug(message string, args ...interface{}) {
  Logger.Debug(message, args...)
}

/*
This will wrap error logging (where there is no exit)
and will handle based on the machine output setting
*/
func LogError(message string, args ...interface{}) {
  if !machineOutput {
    Logger.Error(message)
  }
}

/*
Type to output error json in case of 
machine output being set
*/
type errorOutput struct {
  ExitCode int                `json:"exitCode"`
  ErrorMessage string         `json:"errorMessage"`
}

/*
Get output json for errors 
*/
func (e errorOutput) getOutputJson() (string, int) {
  jsonBytes, err := json.Marshal(e)
  if err != nil {
    return "{\"exitCode\": 100, \"errorMessage\": \"Error marshaling machine output\"}", 100   
  }
  return string(jsonBytes), 0
}

/*
This will wrap error logging with an exit and setting 
exit code with it and will act based on machine output 
setting
*/
func LogErrorExit(meesage string, exitCode int, err error) {
  var machineErrorOutput errorOutput

  if !machineOutput {
    Logger.Error(meesage, "error", err)
    os.Exit(exitCode)
  } else {
    machineErrorOutput.ExitCode = exitCode
    machineErrorOutput.ErrorMessage = fmt.Sprintf("%s: %v", meesage, err)
    output, _ := machineErrorOutput.getOutputJson()
    fmt.Println(output)
    os.Exit(exitCode)
  }
}
