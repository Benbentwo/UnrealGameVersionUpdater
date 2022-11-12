package version

import (
	"github.com/blang/semver"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test(t *testing.T) {
	tests := []struct {
		name          string
		version       string
		wantedVersion string
	}{
		{"GetVersion", "1.2.1", "1.2.1"},
		{"GetVersion", "", "0.0.0"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			Map["version"] = test.version
			result := GetVersion()
			assert.Equal(t, test.wantedVersion, result)
		})
	}
}

func TestGetVersion(t *testing.T) {
	Map["version"] = "1.2.1"
	result := GetVersion()
	assert.Equal(t, "1.2.1", result)
}

func TestGetSemverVersisonWithStandardVersion(t *testing.T) {
	Map["version"] = "1.2.1"
	result, err := GetSemverVersion()
	expectedResult := semver.Version{Major: 1, Minor: 2, Patch: 1}
	assert.NoError(t, err, "GetSemverVersion should exit without failure")
	assert.Exactly(t, expectedResult, result)
}

func TestGetSemverVersisonWithNonStandardVersion(t *testing.T) {
	Map["version"] = "1.3.153-dev+7a8285f4"
	result, err := GetSemverVersion()

	prVersions := []semver.PRVersion{
		semver.PRVersion{VersionStr: "dev"},
	}
	builds := []string{"7a8285f4"}
	expectedResult := semver.Version{Major: 1, Minor: 3, Patch: 153, Pre: prVersions, Build: builds}
	assert.NoError(t, err, "GetSemverVersion should exit without failure")
	assert.Exactly(t, expectedResult, result)
}

func TestGetSemverVersisonErr(t *testing.T) {
	Map["version"] = "abc"
	result, err := GetSemverVersion()
	expectedResult := semver.Version{Major: 0, Minor: 0, Patch: 0, Pre: nil, Build: nil}
	assert.EqualError(t, err, "failed to parse version abc: No Major.Minor.Patch elements found")
	assert.Equal(t, expectedResult, result)
}

func TestVersionStringDefault(t *testing.T) {
	Map["version"] = "1.2.3"
	result := VersionStringDefault("a")
	assert.Equal(t, "1.2.3", result)

	Map["version"] = "abc"
	result = VersionStringDefault("a")
	assert.Equal(t, "a", result)
}
