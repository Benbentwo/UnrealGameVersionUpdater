package common

import (
	"fmt"
	"github.com/Benbentwo/UnrealGameVersionUpdater/pkg/common/log"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"gopkg.in/AlecAivazis/survey.v1/terminal"
	"io"
	"net/url"
	"os"
	"testing"
)

func TestCheckErr(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
	}{
		{"No Error", args{nil}},
	}

	assert.NotPanics(t, func() {
		CheckErr(tests[0].args.err)
	})
}
func Test_checkErr(t *testing.T) {
	ret := func(s string, i int) {
		log.Logger().Infof("%d:%s", i, s)
	}

	type args struct {
		err       error
		handleErr func(string, int)
	}
	tests := []struct {
		name     string
		args     args
		expected string
	}{
		{"No Error", args{nil, ret}, ""},
		{"Fmt Exit", args{fmt.Errorf("exit"), ret}, "INFO: 1:error: exit\n"},
		{"ErrExit", args{ErrExit, ret}, "INFO: 1:\n"},
		{"Spaghetti", args{fmt.Errorf("spaghetti"), ret}, "INFO: 1:error: spaghetti\n"},
		{"E Tacos", args{fmt.Errorf("error: tacos"), ret}, "INFO: 1:error: tacos\n"},
		{"EE Tacos", args{fmt.Errorf("error: error: tacos"), ret}, "INFO: 1:error: error: tacos\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logs := log.CaptureOutput(func() {
				checkErr(tt.args.err, tt.args.handleErr)
			})
			assert.Equal(t, tt.expected, logs)
		})
	}
}

func TestCommonOptions_AddBaseFlags(t *testing.T) {
	type fields struct {
		Cmd       *cobra.Command
		Args      []string
		BatchMode bool
		Verbose   bool
		Quiet     bool
		In        terminal.FileReader
		Out       terminal.FileWriter
		Err       io.Writer
	}
	type args struct {
		cmd *cobra.Command
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"No Opts", fields{}, args{}},
		{"quiet", fields{Quiet: false}, args{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &CommonOptions{
				Cmd:       tt.fields.Cmd,
				Args:      tt.fields.Args,
				BatchMode: tt.fields.BatchMode,
				Verbose:   tt.fields.Verbose,
				Quiet:     tt.fields.Quiet,
				In:        tt.fields.In,
				Out:       tt.fields.Out,
				Err:       tt.fields.Err,
			}
			logCommand := &cobra.Command{
				Use:   "test",
				Short: "dummy test",
			}

			o.AddBaseFlags(logCommand)

			batch, err := logCommand.PersistentFlags().GetBool(OptionBatchMode)
			assert.NoError(t, err)
			assert.Equal(t, batch, tt.fields.BatchMode)

			verbose, err := logCommand.PersistentFlags().GetBool(OptionVerbose)
			assert.NoError(t, err)
			assert.Equal(t, verbose, tt.fields.Verbose)

			quiet, err := logCommand.PersistentFlags().GetBool(OptionQuiet)
			assert.NoError(t, err)
			assert.Equal(t, quiet, tt.fields.Quiet)
		})
	}
}

func TestFatal(t *testing.T) {
	t.Cleanup(func() {
		osExit = os.Exit
	})
	osExit = func(code int) {
		log.Logger().Errorf("supposed to exit with %d", code)
	}
	type args struct {
		message string
		code    int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"No Message", args{message: "", code: 1}, "ERROR: supposed to exit with 1\n"},
		{"Message", args{message: "Blah BLah BLAH", code: 2}, "ERROR: supposed to exit with 2\n"},
		{"Message with \n suffix", args{message: "Blah BLah BLAH\n", code: 3}, "ERROR: supposed to exit with 3\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logs := log.CaptureOutput(func() {
				Fatal(tt.args.message, tt.args.code)
			})
			assert.Equal(t, tt.want, logs)
		})
	}

}

func TestStandardErrorMessage(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 bool
	}{
		{"No Error", args{nil}, "", false},
		{"Network Err, fake", args{&url.Error{Op: "OP", URL: "URL", Err: fmt.Errorf("fake")}}, "Unable to connect to the server: fake", true},
		{"Network Err, conn refused", args{&url.Error{Op: "OP", URL: "http://localhost", Err: fmt.Errorf("connection refused")}}, "The connection to the server localhost was refused - did you specify the right host or port?", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := StandardErrorMessage(tt.args.err)
			if got != tt.want {
				t.Errorf("StandardErrorMessage() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("StandardErrorMessage() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
