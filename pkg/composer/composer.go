// Package composer contains functions and structures
// for working with composer.json.
package composer

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// Config is a structure that stores all the required fields
// from composer.json.
type Config struct {
	Name        string            `json:"name"`
	Type        string            `json:"type"`
	Require     map[string]string `json:"require"`
	RequireDev  map[string]string `json:"require-dev"`
	Reps        []*ConfigRepo     `json:"repositories"`
	Autoload    Autoload          `json:"autoload"`
	AutoloadDev Autoload          `json:"autoload-dev"`

	// Path to the config.
	Path string
	// RootDir is a dir with config.
	RootDir string

	Checks []func(*Config) *ConfigError
}

// Autoload structure stores a mapping to namespaces
// and their actual folders.
//
// See root.handleNamespace function
type Autoload struct {
	Psr4  map[string]string `json:"psr-4"`
	Files []string          `json:"files"`
}

// Psr4PathForNamespace for the passed namespace looks for the path
// in the current autoload psr-4 field.
//
// The search is not performed verbatim, it is enough
// that the namespace is a prefix of one of the psr-4 map keys.
func (a *Autoload) Psr4PathForNamespace(name string) (string, bool) {
	// Since names in psr-4 always end with a slash, we need to add a
	// slash to the namespace name to properly handle the case when
	// the namespace name is equal to the name in psr-4.
	// Example:
	//   psr-4: "My\\Core\\"
	//   name:  "My\\Core"
	name = name + `\`

	// Since we are looking for a prefix, there may be a situation
	// where there are several matches.
	// Example:
	//   name:   My/Core/Utils
	//   match1: My/
	//   match2: My/Core/
	//
	// So, we have to choose the largest prefix found in order to work correctly.
	var psrNameForFound string
	var foundPath string

	for psrName, psrPath := range a.Psr4 {
		if strings.HasPrefix(name, psrName) {
			if len(psrName) > len(psrNameForFound) {
				psrNameForFound = psrName
				foundPath = psrPath
			}
		}
	}

	if foundPath == "" {
		return "", false
	}

	return foundPath, true
}

// Psr4PathForNamespace for the passed namespace looks for the path
// in the autoload.psr-4 and autoload-dev.psr-4 fields.
//
// Returns the found path starting with the folder where
// the microservice or package is located.
//
// See Autoload.Psr4PathForNamespace
func (c *Config) Psr4PathForNamespace(name string) (string, bool) {
	// We need to add a folder to the resulting path to
	// avoid triggers when there is a folder with the
	// same name in the path.
	// Example:
	//   resultPath: src
	//   path:       core/tests/some/src/
	dir := filepath.Base(c.RootDir)

	path, contains := c.Autoload.Psr4PathForNamespace(name)
	if contains {
		return dir + "/" + path, true
	}

	path, containsInDev := c.AutoloadDev.Psr4PathForNamespace(name)
	if containsInDev {
		return dir + "/" + path, true
	}

	return "", false
}

// ConfigRepo is a structure for storing dependencies
// from composer.json.
type ConfigRepo struct {
	Type     string `json:"type"`
	Url      string `json:"url"`
	Resolved bool
}

// NewConfigFromFile returns new config from file.
//
// If the file does not exist or contains invalid json an error will be returned.
func NewConfigFromFile(path string) (*Config, *ConfigErrors) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return &Config{}, NewConfigErrors(&ConfigError{
			Msg:      err.Error(),
			Critical: true,
		})
	}
	return NewConfigFromData(data, path)
}

// NewConfigFromData returns new config from data.
//
// If data contains invalid json an error will be returned.
func NewConfigFromData(data []byte, configPath string) (*Config, *ConfigErrors) {
	var config Config
	err := json.Unmarshal(data, &config)
	if err != nil {
		return &Config{}, NewConfigErrors(&ConfigError{
			Msg:      err.Error(),
			Critical: true,
		})
	}
	absPath, _ := filepath.Abs(configPath)
	root := filepath.Dir(absPath)

	config.Path = absPath

	config.RootDir = root
	return &config, nil
}

// AddCheck adds custom check for config.
func (c *Config) AddCheck(check func(*Config) *ConfigError) {
	c.Checks = append(c.Checks, check)
}

// CheckConfig checks the config against the rules.
//
// See Config.AddCheck
func (c *Config) CheckConfig() *ConfigErrors {
	errors := &ConfigErrors{
		Config: c,
	}

	for _, check := range c.Checks {
		if err := check(c); err != nil {
			errors.Add(err)
		}
	}

	if errors.Len() == 0 {
		return nil
	}

	return errors
}

// ResolveUrl resolves the path for the dependency relative to the passed path.
func (c *ConfigRepo) ResolveUrl(path string) *ConfigRepo {
	if c.Resolved {
		return c
	}

	// At the moment it is not clear how to handle non-local
	// inclusions (for example via vcs).
	if c.Type != "path" {
		return c
	}

	c.Url = filepath.Clean(filepath.Join(path, c.Url))

	// In order to correctly handle paths in unix-like systems and in windows,
	// we need to bring all slashes to the form as in unix.
	c.Url = filepath.ToSlash(c.Url)

	c.Resolved = true
	return c
}
