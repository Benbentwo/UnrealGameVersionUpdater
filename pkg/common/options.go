package common

import (
	"fmt"
	"github.com/Benbentwo/go-bin-generic/pkg/common/log"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"gopkg.in/AlecAivazis/survey.v1/terminal"
	"io"
	"net/url"
	"os"
	"strconv"

	"strings"
)

const (
	OptionBatchMode = "batch-mode"
	OptionVerbose   = "verbose"
	OptionQuiet     = "quiet" // sets to 	warn 	level

)

type CommonOptions struct {
	Cmd       *cobra.Command
	Args      []string
	BatchMode bool
	Verbose   bool
	Quiet     bool
	In        terminal.FileReader
	Out       terminal.FileWriter
	Err       io.Writer
}

// AddBaseFlags adds the base flags for all commands
func (o *CommonOptions) AddBaseFlags(cmd *cobra.Command) {
	defaultBatchMode := false
	if os.Getenv("BATCH_MODE") == "true" {
		defaultBatchMode = true
	}
	cmd.PersistentFlags().BoolVarP(&o.BatchMode, OptionBatchMode, "b", defaultBatchMode, "Runs in batch mode without prompting for user input")
	cmd.PersistentFlags().BoolVarP(&o.Verbose, OptionVerbose, "", false, "Enables verbose output")
	cmd.PersistentFlags().BoolVarP(&o.Quiet, OptionQuiet, "q", false, "Enables quiet output")

	o.Cmd = cmd
}

func SetLoggingLevel(cmd *cobra.Command, args []string) {
	verbose, _ := strconv.ParseBool(cmd.Flag(OptionVerbose).Value.String())
	quiet, _ := strconv.ParseBool(cmd.Flag(OptionQuiet).Value.String())
	level := os.Getenv("VGS_LOG_LEVEL")
	if level != "" {
		if verbose {
			log.Logger().Trace("The VGS_LOG_LEVEL environment variable took precedence over the verbose flag")
		}

		err := log.SetLevel(level)
		if err != nil {
			log.Logger().Errorf("Unable to set log level to %s", level)
			level = ""
		}
	}
	if level == "" {
		if verbose {
			_ = log.SetLevel("debug")
		} else if quiet {
			_ = log.SetLevel("warn")
		} else {
			_ = log.SetLevel("info")
		}
	}
}

const (
	defaultErrorExitCode = 1
)

var osExit = os.Exit
var fatalErrHandler = Fatal

// BehaviorOnFatal allows you to override the default behavior when a fatal
// error occurs, which is to call os.Exit(code). You can pass 'panic' as a function
// here if you prefer the panic() over os.Exit(1).
func BehaviorOnFatal(f func(string, int)) {
	fatalErrHandler = f
}

// DefaultBehaviorOnFatal allows you to undo any previous override.  Useful in tests.
func DefaultBehaviorOnFatal() {
	fatalErrHandler = Fatal
}

// Fatal prints the message (if provided) and then exits. If V(2) or greater,
// glog.Logger().Fatal is invoked for extended information.
func Fatal(msg string, code int) {
	if len(msg) > 0 {
		// add newline if needed
		if !strings.HasSuffix(msg, "\n") {
			msg += "\n"
		}
		fmt.Fprint(os.Stderr, msg)
	}
	osExit(code)
}

// This method is generic to the command in use and may be used by non-Kubectl
// commands.
func CheckErr(err error) {
	checkErr(err, fatalErrHandler)
}

// ErrExit may be passed to CheckError to instruct it to output nothing but exit with
// status code 1.
var ErrExit = fmt.Errorf("exit")

// checkErr formats a given error as a string and calls the passed handleErr
// func with that string and an kubectl exit code.
func checkErr(err error, handleErr func(string, int)) {
	switch {
	case err == nil:
		return
	case err == ErrExit:
		handleErr("", defaultErrorExitCode)
		return
	default:
		switch err := err.(type) {
		default: // for any other error type
			msg, ok := StandardErrorMessage(err)
			if !ok {
				msg = err.Error()
				if !strings.HasPrefix(msg, "error: ") {
					msg = fmt.Sprintf("error: %s", msg)
				}
			}
			handleErr(msg, defaultErrorExitCode)
		}
	}
}

// StandardErrorMessage translates common errors into a human readable message, or returns
// false if the error is not one of the recognized types. It may also log extended
// information to glog.
//
// This method is generic to the command in use and may be used by non-Kubectl
// commands.
func StandardErrorMessage(err error) (string, bool) {
	switch t := err.(type) {
	case *url.Error:
		glog.V(4).Infof("Connection error: %s %s: %v", t.Op, t.URL, t.Err)
		switch {
		case strings.Contains(t.Err.Error(), "connection refused"):
			host := t.URL
			if server, err := url.Parse(t.URL); err == nil {
				host = server.Host
			}
			return fmt.Sprintf("The connection to the server %s was refused - did you specify the right host or port?", host), true
		}
		return fmt.Sprintf("Unable to connect to the server: %v", t.Err), true
	}
	return "", false
}
