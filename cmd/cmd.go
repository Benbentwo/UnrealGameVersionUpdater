package cmd

import (
	"github.com/Benbentwo/UnrealGameVersionUpdater/pkg/common"
	"github.com/Benbentwo/UnrealGameVersionUpdater/pkg/version"
	"github.com/spf13/viper"
	"io"
	"strings"

	"github.com/Benbentwo/utils/util"
	"github.com/spf13/cobra"
	"gopkg.in/AlecAivazis/survey.v1/terminal"
)

// Build information. Populated at build-time.
var (
	Binary string
)

func NewMainCmd(in terminal.FileReader, out terminal.FileWriter, err io.Writer, args []string) *cobra.Command {

	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	cmd := &cobra.Command{
		Use:              Binary,
		Short:            "CLI tool",
		Long:             "This CLI tool is designed to help get started with other projects, TODO: CHANGEME",
		PersistentPreRun: common.SetLoggingLevel,
		Run:              runHelp,
	}
	commonOpts := &common.CommonOptions{
		In:  in,
		Out: out,
		Err: err,
	}
	commonOpts.AddBaseFlags(cmd)

	// Section to add commands to:
	cmd.AddCommand(version.NewCmdVersion(commonOpts))

	return cmd
}

func runHelp(cmd *cobra.Command, args []string) {
	err := cmd.Help()
	if err != nil {
		util.Logger().Errorf("Error running help: %s", err)
	}
}
