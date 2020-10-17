package composer

import (
	"fmt"
)

// ConfigError structure describes one error in the config.
//
// If the Critical flag is true, then the analysis
// process will be interrupted.
type ConfigError struct {
	Msg      string
	Critical bool
}

// Error returns a string with error message and critical flag.
func (ce ConfigError) Error() string {
	var res string
	if ce.Critical {
		res += "<critical> "
	}
	res += ce.Msg
	return res
}

// ConfigErrors is a structure for storing all errors in the config.
type ConfigErrors struct {
	Config *Config
	Errors []*ConfigError
}

// NewConfigErrors creates a set of config errors from passed errors.
func NewConfigErrors(errors ...*ConfigError) *ConfigErrors {
	return &ConfigErrors{
		Errors: errors,
	}
}

// Add adds a new error.
func (ce *ConfigErrors) Add(err *ConfigError) {
	ce.Errors = append(ce.Errors, err)
}

// Len returns the number of errors.
func (ce *ConfigErrors) Len() int {
	return len(ce.Errors)
}

// Error returns a string with one error on each line.
func (ce *ConfigErrors) Error() string {
	var res string
	for _, e := range ce.Errors {
		res += fmt.Sprintf("config %s: %s\n", ce.Config.Path, e.Error())
	}
	return res
}
