package gosec

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
)

const (
	// Globals are applicable to all rules and used for general
	// configuration settings for gosec.
	Globals = "global"
)

// Config is used to provide configuration and customization to each of the rules.
type Config map[string]interface{}

// NewConfig initializes a new configuration instance. The configuration data then
// needs to be loaded via c.ReadFrom(strings.NewReader("config data"))
// or from a *os.File.
func NewConfig() Config {
	cfg := make(Config)
	cfg[Globals] = make(map[string]string)
	return cfg
}

// ReadFrom implements the io.ReaderFrom interface. This
// should be used with io.Reader to load configuration from
//file or from string etc.
func (c Config) ReadFrom(r io.Reader) (int64, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return int64(len(data)), err
	}
	if err = json.Unmarshal(data, &c); err != nil {
		return int64(len(data)), err
	}
	return int64(len(data)), nil
}

// WriteTo implements the io.WriteTo interface. This should
// be used to save or print out the configuration information.
func (c Config) WriteTo(w io.Writer) (int64, error) {
	data, err := json.Marshal(c)
	if err != nil {
		return int64(len(data)), err
	}
	return io.Copy(w, bytes.NewReader(data))
}

// Get returns the configuration section for the supplied key
func (c Config) Get(section string) (interface{}, error) {
	settings, found := c[section]
	if !found {
		return nil, fmt.Errorf("Section %s not in configuration", section)
	}
	return settings, nil
}

// Set section in the configuration to specified value
func (c Config) Set(section string, value interface{}) {
	c[section] = value
}

// GetGlobal returns value associated with global configuration option
func (c Config) GetGlobal(option string) (string, error) {
	if globals, ok := c[Globals]; ok {
		if settings, ok := globals.(map[string]string); ok {
			if value, ok := settings[option]; ok {
				return value, nil
			}
			return "", fmt.Errorf("global setting for %s not found", option)
		}
	}
	return "", fmt.Errorf("no global config options found")

}

// SetGlobal associates a value with a global configuration ooption
func (c Config) SetGlobal(option, value string) {
	if globals, ok := c[Globals]; ok {
		if settings, ok := globals.(map[string]string); ok {
			settings[option] = value
		}
	}
}
