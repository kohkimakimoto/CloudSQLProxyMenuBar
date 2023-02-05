package main

import (
	"github.com/BurntSushi/toml"
	"sort"
)

const InitialConfig = `# This is the main configuration file of CloudSQLProxyMenuBar.
# https://github.com/kohkimakimoto/CloudSQLProxyMenuBar

#
# core is the section of CloudSQLProxyMenuBar global config.
#
[core]
# Optional: The path to 'cloud_sql_proxy' command.
# If you don't specify it, CloudSQLProxyMenuBar uses 'cloud_sql_proxy' command in your PATH.
# If you are not familiar with cloud_sql_proxy, please read the document: https://cloud.google.com/sql/docs/mysql/sql-proxy
# cloud_sql_proxy = "/path/to/cloud_sql_proxy"

# Optional: The log file path.
# The default is '$HOME/.cloudsqlproxymenubar/output.log''
# log_file = "/path/to/logfile"

#
# proxies.xxx is the section of the Cloud SQL Proxy settings.
#
# [proxies.yourinstance]
# Optional: The text is displayed on the menu item.
# The default is the same as 'XXX' part of 'proxies.XXX'.
# label = "proxy-to-yourinstance1"

# Required: The command line arguments passed to 'cloud_sql_proxy' command.
# arguments = "-dir=/cloudsql -instances=yourcompany:asia-northeast1:yourinstance -credential_file=xxx.json"

# You can add proxy config multiple times.
# [proxies.yourinstance2]
# ...
`

type Config struct {
	Core    *CoreConfig             `toml:"core"`
	Proxies map[string]*ProxyConfig `toml:"proxies"`
}

func NewConfig() *Config {
	return &Config{
		Core: &CoreConfig{
			LogFile:       "",
			CloudSqlProxy: "cloud_sql_proxy",
		},
		Proxies: map[string]*ProxyConfig{},
	}
}

func (c *Config) SortedProxyKeys() []string {
	ret := make([]string, 0, len(c.Proxies))
	for key := range c.Proxies {
		ret = append(ret, key)
	}
	sort.Strings(ret)
	return ret
}

type CoreConfig struct {
	LogFile       string `toml:"log_file"`
	CloudSqlProxy string `toml:"cloud_sql_proxy"`
}

type ProxyConfig struct {
	Name      string `toml:"-"`
	Label     string `toml:"label"`
	Arguments string `toml:"arguments"`
}

func (c *ProxyConfig) LabelOrName() string {
	if c.Label != "" {
		return c.Label
	}
	return c.Name
}

func (c *Config) LoadFromFile(path string) error {
	_, err := toml.DecodeFile(path, c)
	if err != nil {
		return err
	}

	// set proxy name
	for key, proxy := range c.Proxies {
		proxy.Name = key
	}

	return nil
}
