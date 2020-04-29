package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersionedConfig_GetData(t *testing.T) {
	vc := VersionedConfig{
		data: []byte{0x00, 0x01, 0x02},
	}
	assert.Equal(t, []byte{0x00, 0x01, 0x02}, vc.GetData())
}

func TestVersionedConfig_GetData_nil(t *testing.T) {
	var vc *VersionedConfig
	assert.Equal(t, []byte(nil), vc.GetData())
}

func TestVersionedConfig_GetData_nil2(t *testing.T) {
	var vc VersionedConfig
	assert.Equal(t, []byte(nil), vc.GetData())
}

func TestVersionedConfig_GetPath(t *testing.T) {
	vc := VersionedConfig{
		path: "test/path",
	}
	assert.Equal(t, "test/path", vc.GetPath())
}

func TestVersionedConfig_GetPath_nil(t *testing.T) {
	var vc *VersionedConfig
	assert.Equal(t, "", vc.GetPath())
}

func TestVersionedConfig_GetPath_nil2(t *testing.T) {
	var vc VersionedConfig
	assert.Equal(t, "", vc.GetPath())
}

func TestVersionedConfig_GetVersion(t *testing.T) {
	s := "v1"
	vc := VersionedConfig{
		Version: &s,
	}
	assert.Equal(t, "v1", vc.GetVersion())
}

func TestVersionedConfig_GetVersion_nil(t *testing.T) {
	var vc *VersionedConfig
	assert.Equal(t, "", vc.GetVersion())
}

func TestVersionedConfig_GetVersion_nil2(t *testing.T) {
	var vc VersionedConfig
	assert.Equal(t, "", vc.GetVersion())
}

func TestVersionedConfig_GetVersion_noVersion(t *testing.T) {
	vc := VersionedConfig{}
	assert.Equal(t, "", vc.GetVersion())
}

func TestLoad(t *testing.T) {
	vc, err := Load("testdata/config.yaml")
	assert.NoError(t, err)

	assert.Equal(t, "v1", vc.GetVersion())
	assert.Equal(t, "testdata/config.yaml", vc.GetPath())
	assert.Equal(t, "version: v1\nchart:\n  name: chart\n  repo: github.com/edaniszewski/charts-test\n", string(vc.GetData()))
}

func TestLoad_ConfigPathNotExist(t *testing.T) {
	// Load will attempt to append the default file name to the directory, which
	// does not exist in the testdata path.
	_, err := Load("testdata/")
	assert.Error(t, err)
	assert.Equal(t, ErrNoConfig, err)
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := Load("testdata/not_found.yaml")
	assert.EqualError(t, err, "stat testdata/not_found.yaml: no such file or directory")
}

func TestLoad_InvalidYAML(t *testing.T) {
	_, err := Load("testdata/invalid.yaml")
	assert.EqualError(t, err, "yaml: line 1: did not find expected node content")
	assert.Error(t, err)
}

func TestLoad_NoVersion(t *testing.T) {
	_, err := Load("testdata/config_no_version.yaml")
	assert.Error(t, err)
	assert.Equal(t, ErrNoConfigVersion, err)
}
