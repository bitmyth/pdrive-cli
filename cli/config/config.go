package config

import (
	baseConfig "github.com/bitmyth/pdrive-cli/cli/pkg/config"
	"os"
	"path/filepath"
)

const (
	hosts   = "hosts"
	aliases = "aliases"
)

type Config interface {
	AuthToken(string) (string, string)
	Get(string, string) (string, error)
	GetOrDefault(string, string) (string, error)
	Set(string, string, string)
	UnsetHost(string)
	Hosts() ([]string, string)
	DefaultHost() (string, string)
	DefaultScheme() (string, string)
	Write() error
}

func NewConfig() (Config, error) {
	c, err := baseConfig.Read()
	if err != nil {
		return nil, err
	}
	return &cfg{c}, nil
}

// Implements Config interface
type cfg struct {
	cfg *baseConfig.Config
}

func (c *cfg) AuthToken(host string) (string, string) {
	return baseConfig.TokenForHost(host)
}

func (c *cfg) Get(hostname, key string) (string, error) {
	if hostname != "" {
		val, err := c.cfg.Get([]string{hosts, hostname, key})
		if err == nil {
			return val, err
		}
	}

	return c.cfg.Get([]string{key})
}

func (c *cfg) GetOrDefault(hostname, key string) (string, error) {
	var val string
	var err error
	if hostname != "" {
		val, err = c.cfg.Get([]string{hosts, hostname, key})
		if err == nil {
			return val, err
		}
	}

	val, err = c.cfg.Get([]string{key})
	if err == nil {
		return val, err
	}

	if defaultExists(key) {
		return defaultFor(key), nil
	}

	return val, err
}

func (c *cfg) Set(hostname, key, value string) {
	if hostname == "" {
		c.cfg.Set([]string{key}, value)
		return
	}
	c.cfg.Set([]string{hosts, hostname, key}, value)
}

func (c *cfg) UnsetHost(hostname string) {
	if hostname == "" {
		return
	}
	_ = c.cfg.Remove([]string{hosts, hostname})
}

func (c *cfg) Hosts() ([]string, string) {
	return baseConfig.Hosts()
}

func (c *cfg) DefaultScheme() (string, string) {
	return baseConfig.DefaultScheme()
}

func (c *cfg) DefaultHost() (string, string) {
	return baseConfig.DefaultHost()
}

func (c *cfg) Aliases() *AliasConfig {
	return &AliasConfig{cfg: c.cfg}
}

func (c *cfg) Write() error {
	return baseConfig.Write(c.cfg)
}

func defaultFor(key string) string {
	for _, co := range configOptions {
		if co.Key == key {
			return co.DefaultValue
		}
	}
	return ""
}

func defaultExists(key string) bool {
	for _, co := range configOptions {
		if co.Key == key {
			return true
		}
	}
	return false
}

type AliasConfig struct {
	cfg *baseConfig.Config
}

func (a *AliasConfig) Get(alias string) (string, error) {
	return a.cfg.Get([]string{aliases, alias})
}

func (a *AliasConfig) Add(alias, expansion string) {
	a.cfg.Set([]string{aliases, alias}, expansion)
}

func (a *AliasConfig) Delete(alias string) error {
	return a.cfg.Remove([]string{aliases, alias})
}

func (a *AliasConfig) All() map[string]string {
	out := map[string]string{}
	keys, err := a.cfg.Keys([]string{aliases})
	if err != nil {
		return out
	}
	for _, key := range keys {
		val, _ := a.cfg.Get([]string{aliases, key})
		out[key] = val
	}
	return out
}

type ConfigOption struct {
	Key           string
	Description   string
	DefaultValue  string
	AllowedValues []string
}

var configOptions = []ConfigOption{
	{
		Key:          "default_host",
		Description:  "default host",
		DefaultValue: "pdrive.danyuan.ink",
	},
	{
		Key:          "env",
		Description:  "env",
		DefaultValue: "production",
	},
	{
		Key:          "editor",
		Description:  "the text editor program to use for authoring text",
		DefaultValue: "",
	},
	{
		Key:           "prompt",
		Description:   "toggle interactive prompting in the terminal",
		DefaultValue:  "enabled",
		AllowedValues: []string{"enabled", "disabled"},
	},
	{
		Key:          "pager",
		Description:  "the terminal pager program to send standard output to",
		DefaultValue: "",
	},
	{
		Key:          "http_unix_socket",
		Description:  "the path to a Unix socket through which to make an HTTP connection",
		DefaultValue: "",
	},
	{
		Key:          "browser",
		Description:  "the web browser to use for opening URLs",
		DefaultValue: "",
	},
}

func ConfigOptions() []ConfigOption {
	return configOptions
}

func HomeDirPath(subdir string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	newPath := filepath.Join(homeDir, subdir)
	return newPath, nil
}

func DataDir() string {
	return baseConfig.DataDir()
}

func ConfigDir() string {
	return baseConfig.ConfigDir()
}
