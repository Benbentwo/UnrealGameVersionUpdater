package version

import (
	"fmt"
	"github.com/Benbentwo/go-bin-generic/pkg/github"
	"github.com/Benbentwo/utils/util"
	"github.com/blang/semver"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const ( // iota is reset to 0
	MAC_OS      = iota
	LINUX_ARM   = iota
	LINUX_AMD64 = iota
	WIN_AMD64   = iota
	WIN_i386    = iota
)

// GetLatestVersion returns latest version
func (o *VersionOptions) GetLatestVersion() (semver.Version, string, error) {

	// if runtime.GOOS == "darwin" {
	// 	o.OS = MAC_OS
	// 	util.Logger().Debugf("Using MacOS")
	// }
	return github_helpers.GetLatestVersionFromGitHub(Org, Repo)
}

// InstallBin installs repo's cli
func (o *VersionOptions) InstallBin(upgrade bool, prefix string, version string) error {
	util.Logger().Debugf("installing "+Repo+" %s", version)
	if runtime.GOOS == "darwin" {
	}
	binDir, err := BinLocation()
	if err != nil {
		return err
	}
	// Check for binary in non standard path and install there instead if found...
	nonStandardBinDir, err := BinaryLocation()
	if err == nil && binDir != nonStandardBinDir {
		binDir = nonStandardBinDir
	}
	binary := Binary
	fileName := binary
	if !upgrade {
		f, flag, err := ShouldInstallBinary(binary)
		if err != nil || !flag {
			return err
		}
		fileName = f
	}
	if version == "" {
		latestVersion, latestPrefix, err := github_helpers.GetLatestVersionFromGitHub(Org, Repo)
		if err != nil {
			return err
		}
		version = fmt.Sprintf("%s", latestVersion)
		prefix = fmt.Sprintf("%s", latestPrefix)
	}

	// NOTE the below is pretty similar to Jenkins-x Version, the primary difference is that This repository is uploading
	// Basic binaries and exectuables, not tar.gz versions. so extension needs to include .exe or be nothing
	extension := ""
	if runtime.GOOS == "windows" {
		extension = ".exe"
	}

	protocol := "https://"
	if strings.HasPrefix(GitServer, "http") {
		// Project configured to manually set protocol.
		protocol = ""
	}

	// Set in the makefile
	BinaryDownloadBaseURL := strings.Join([]string{GitServer, Org, Repo, "releases", "download", prefix}, "/")
	// BinaryDownloadBaseURL := "https://github.com/Benbentwo/go-bin-generic/releases/download/v1.0.0/
	// 							 https://github.com/Benbentwo/go-bin-generic/releases/download/v1.0.0/go-bin-generic-windows-amd64.zip"
	clientURL := fmt.Sprintf("%s%s%s/"+binary+"-%s-%s%s", protocol, BinaryDownloadBaseURL, version, runtime.GOOS, runtime.GOARCH, extension)
	fullPath := filepath.Join(binDir, fileName)
	if runtime.GOOS == "windows" {
		fullPath += ".exe"
	}
	tmpArchiveFile := fullPath + ".tmp"
	err = DownloadFile(clientURL, tmpArchiveFile)
	if err != nil {
		return err
	}

	if runtime.GOOS != "windows" {
		err = os.Rename(tmpArchiveFile, fullPath)
		if err != nil {
			return err
		}
	} else { // windows
		// A standard remove and rename (or overwrite) will not work as the file will be locked as windows is running it
		// the trick is to rename to a tempfile :-o
		// this will leave old files around but well at least it updates.
		// we could schedule the file for cleanup at next boot but....
		// HKLM\System\CurrentControlSet\Control\Session Manager\PendingFileRenameOperations
		err = os.Rename(filepath.Join(binDir, binary+".exe"), filepath.Join(binDir, binary+".exe.deleteme"))
		// if we can not rename it this i pretty fatal as we won;t be able to overwrite either
		if err != nil {
			return err
		}
		// Copy over the new binary
		err = os.Rename(tmpArchiveFile, fullPath)
		if err != nil {
			return err
		}
	}
	util.Logger().Infof(util.ColorDebug(Repo)+" cli has been installed into %s", util.ColorInfo(fullPath))
	return os.Chmod(fullPath, 0755)
}

func BinLocation() (string, error) {
	h := util.HomeDir()
	path := filepath.Join(h, "bin")
	err := os.MkdirAll(path, util.DefaultWritePermissions)
	if err != nil {
		return "", err
	}
	return path, nil
}

// BinaryLocation Returns the path to the currently installed binary.
func BinaryLocation() (string, error) {
	return binaryLocation(os.Executable)
}

func binaryLocation(osExecutable func() (string, error)) (string, error) {
	processBinary, err := osExecutable()
	if err != nil {
		util.Logger().Debugf("processBinary error %s", err)
		return processBinary, err
	}
	util.Logger().Debugf("processBinary %s", processBinary)
	// make it absolute
	processBinary, err = filepath.Abs(processBinary)
	if err != nil {
		util.Logger().Debugf("processBinary error %s", err)
		return processBinary, err
	}
	util.Logger().Debugf("processBinary %s", processBinary)

	// if the process was started form a symlink go and get the absolute location.
	processBinary, err = filepath.EvalSymlinks(processBinary)
	if err != nil {
		util.Logger().Debugf("processBinary error %s", err)
		return processBinary, err
	}

	util.Logger().Debugf("processBinary %s", processBinary)
	path := filepath.Dir(processBinary)
	util.Logger().Debugf("dir from '%s' is '%s'", processBinary, path)
	return path, nil
}

// ShouldInstallBinary checks if the given binary should be installed
func ShouldInstallBinary(name string) (fileName string, download bool, err error) {
	fileName = BinaryWithExtension(name)
	download = false
	pgmPath, err := exec.LookPath(fileName)
	if err == nil {
		util.Logger().Debugf("%s is already available on your PATH at %s", util.ColorInfo(fileName), util.ColorInfo(pgmPath))
		return
	}

	binDir, err := BinLocation()
	if err != nil {
		return
	}

	// lets see if its been installed but just is not on the PATH
	exists, err := util.FileExists(filepath.Join(binDir, fileName))
	if err != nil {
		return
	}
	if exists {
		util.Logger().Debugf("Please add %s to your PATH", util.ColorInfo(binDir))
		return
	}
	download = true
	return
}

func BinaryWithExtension(binary string) string {
	if runtime.GOOS == "windows" {
		if binary == "gcloud" {
			return binary + ".cmd"
		}
		return binary + ".exe"
	}
	return binary
}

// DownloadFile downloads binary content of given URL into local filesystem.
func DownloadFile(clientURL string, fullPath string) error {
	util.Logger().Infof("Downloading %s to %s...", util.ColorInfo(clientURL), util.ColorInfo(fullPath))
	err := DownloadFileFromUrl(fullPath, clientURL)
	if err != nil {
		return fmt.Errorf("Unable to download file %s from %s due to: %v", fullPath, clientURL, err)
	}
	util.Logger().Infof("Downloaded %s", util.ColorInfo(fullPath))
	return nil
}

// Download a file from the given URL
func DownloadFileFromUrl(filepath string, url string) (err error) {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := GetClientWithTimeout(time.Minute * 5).Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("download of %s failed with return code %d", url, resp.StatusCode)
		return err
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	// make it executable
	os.Chmod(filepath, 0755)
	if err != nil {
		return err
	}
	return nil
}

// GetClientWithTimeout returns a client with default transport and user specified timeout
func GetClientWithTimeout(duration time.Duration) *http.Client {
	client := http.Client{
		Timeout: duration,
	}
	return &client
}
