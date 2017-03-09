package models

import (
	"os"

	"github.com/alex-oleshkevich/vhost/utils"
	"gopkg.in/yaml.v2"
)

// LockFilename path
const LockFilename = "vhost.lock"

// Lock file
type Lock struct {
	DbName        string `yaml:"db_name"`
	DbType        string `yaml:"db_type"`
	VhostLinkPath string `yaml:"vhost-link-path"`
}

// Write lockfile
func (l *Lock) Write() error {
	file, err := os.Create(LockFilename)
	if err != nil {
		return err
	}

	contents, err := yaml.Marshal(l)
	if err != nil {
		return err
	}
	_, err = file.Write(contents)
	return err
}

// Read lockfile
func (l *Lock) Read() error {
	file, err := os.Open(LockFilename)
	if err != nil {
		return err
	}

	var contents []byte
	_, err = file.Read(contents)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(contents, l)
}

// Exists tests if lock file exists
func (l *Lock) Exists() bool {
	return utils.FileExists(LockFilename)
}
