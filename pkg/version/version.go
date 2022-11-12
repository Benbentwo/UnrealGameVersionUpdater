package version

import (
	"fmt"
	"github.com/Benbentwo/go-bin-generic/pkg/common"
	"github.com/Benbentwo/utils/util"
	"github.com/blang/semver"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/AlecAivazis/survey.v1"
)

type VersionOptions struct {
	*common.CommonOptions
	OS int
}

func NewCmdVersion(commonOpts *common.CommonOptions) *cobra.Command {
	options := &VersionOptions{
		CommonOptions: commonOpts,
	}

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version information",
		Run: func(cmd *cobra.Command, args []string) {
			options.Cmd = cmd
			options.Args = args
			err := options.Run()
			common.CheckErr(err)
		},
	}
	return cmd
}

func (o *VersionOptions) Run() error {
	currentVersion, err := GetSemverVersion()
	if err != nil {
		return errors.Wrapf(err, "getting semver version: %s\n%s", err, util.ColorWarning("is this a dev build?"))
	}

	util.Logger().Infof("Version: %s", util.ColorInfo(currentVersion))
	err = o.upgradeIfNeeded(currentVersion)
	return err
}

func (o *VersionOptions) upgradeIfNeeded(currentVersion semver.Version) error {
	newVersion, _, err := o.GetLatestVersion()
	if err != nil {
		return errors.Wrap(err, "getting latest version")
	}
	if currentVersion.LT(newVersion) {
		message := fmt.Sprintf("Would you like to upgrade to the %s version?", util.ColorInfo(newVersion))
		answer := true
		prompt := &survey.Confirm{
			Message: message,
			Default: answer,
			Help:    "This will fetch the latest binary and update your local with it",
		}
		surveyOpts := survey.WithStdio(o.In, o.Out, o.Err)
		err := survey.AskOne(prompt, &answer, nil, surveyOpts)
		if err != nil {
			return err
		}

		if answer {
			candidateInstallVersion, prefix, err := o.candidateInstallVersion()
			if err != nil {
				return err
			}

			if o.needsUpgrade(currentVersion, candidateInstallVersion) {
				shouldUpgrade, err := o.ShouldUpdate(candidateInstallVersion)
				if err != nil {
					return errors.Wrap(err, "failed to determine if we should upgrade")
				}
				if shouldUpgrade {
					return o.InstallBin(true, prefix, candidateInstallVersion.String())
				}
			}
		}
	}
	return nil
}

func (o *VersionOptions) needsUpgrade(currentVersion semver.Version, latestVersion semver.Version) bool {
	if latestVersion.EQ(currentVersion) {
		util.Logger().Infof("You are already on the latest version of "+Repo+" %s", util.ColorInfo(currentVersion.String()))
		return false
	}
	return true
}

func (o *VersionOptions) candidateInstallVersion() (semver.Version, string, error) {
	latestVersion, prefix, err := o.GetLatestVersion()
	if err != nil {
		return semver.Version{}, prefix, errors.Wrap(err, "failed to determine version of latest jx release")
	}
	return latestVersion, prefix, nil
}

// ShouldUpdate checks if CLI version should be updated
func (o *VersionOptions) ShouldUpdate(newVersion semver.Version) (bool, error) {
	util.Logger().Debugf("Checking if should upgrade %s", newVersion)
	currentVersion, err := GetSemverVersion()
	if err != nil {
		return false, err
	}

	if newVersion.GT(currentVersion) {
		// Do not ask to update if we are using a dev build...
		for _, x := range currentVersion.Pre {
			if x.VersionStr == "dev" {
				util.Logger().Debugf("Ignoring possible update as it appears you are using a dev build - %s", currentVersion)
				return false, nil
			}
		}
		return true, nil
	}
	return false, nil
}
