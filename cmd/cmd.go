package cmd

import (
	"github.com/Benbentwo/UnrealGameVersionUpdater/pkg/common"
	"github.com/Benbentwo/UnrealGameVersionUpdater/pkg/common/log"
	"github.com/spf13/viper"
	"io"
	"path"
	"strings"

	"github.com/ryanuber/go-glob"
	"github.com/spf13/cobra"
	"gopkg.in/AlecAivazis/survey.v1/terminal"
	"gopkg.in/ini.v1"
	"io/ioutil"
)

// Build information. Populated at build-time.
var (
	Binary            string
	SectionHeader     = "/Script/EngineSettings.GeneralProjectSettings"
	ProjectVersionKey = "ProjectVersion"
)

type VersionUpdaterOptions struct {
	*common.CommonOptions
	IsProject       bool
	IsPlugin        bool
	ConfigDirectory string
}

func NewMainCmd(in terminal.FileReader, out terminal.FileWriter, err io.Writer, args []string) *cobra.Command {

	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	cmd := &cobra.Command{
		Use:              Binary,
		Short:            "CLI tool to update the version of an unreal project",
		Long:             "CLI tool to update the version of an unreal project",
		PersistentPreRun: common.SetLoggingLevel,
		Run:              run,
		Args:             cobra.ExactArgs(1),
	}
	commonOpts := &common.CommonOptions{
		In:  in,
		Out: out,
		Err: err,
	}
	commonOpts.AddBaseFlags(cmd)
	cmd.Flags().BoolP("project", "p", true, "Is the version being updated a project?")
	cmd.Flags().BoolP("plugin", "l", false, "Is the version being updated a project?")
	cmd.Flags().StringP("config", "c", "Config", "Folder where the ini file to be updated live.")

	return cmd
}

func run(cmd *cobra.Command, args []string) {
	version := args[0]
	log.Logger().Debugf("Setting Version to %s", version)

	configDir, _ := cmd.Flags().GetString("config")

	files, err := ioutil.ReadDir(configDir)
	if err != nil {
		log.Logger().Fatalln(err)
	}

	FileFoundIn := ""

	for _, file := range files {
		foundVersion := ""
		filePath := path.Join(configDir, file.Name())
		if glob.Glob("*.ini", file.Name()) {
			cfg, err := ini.Load(filePath)
			if err != nil {
				log.Logger().Fatalf("Failed to load ini file: %s: %s", filePath, err)
			}
			foundVersion = cfg.Section(SectionHeader).Key(ProjectVersionKey).String()
			log.Logger().Debugf("Found Version:\t%s\t%s", file.Name(), foundVersion)
			if foundVersion != "" {
				FileFoundIn = filePath
				break
			}
		}
	}

	if FileFoundIn == "" {
		log.Logger().Fatalln("Could not find a current version in any *.ini files. please add a current version.")
		log.Logger().Fatalln("Haven't implemented adding in a version if not already found.")
	}

	cfg, err := ini.Load(FileFoundIn)
	cfg.Section(SectionHeader).Key(ProjectVersionKey).SetValue(version)
	cfg.SaveTo(FileFoundIn)

}
