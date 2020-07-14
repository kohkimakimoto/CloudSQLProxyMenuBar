package app

import (
	"github.com/BurntSushi/toml"
	"sort"
)

const InitialConfig = `# see https://github.com/kohkimakimoto/CloudSQLProxyMenuBar
[core]
# cloud_sql_proxy = "/usr/local/bin/cloud_sql_proxy"

# [proxies.your-instance]
# label = "production-db"
# options = "-dir=/cloudsql -instances=yourcompany:asia-northeast1:yourinstance -credential_file=xxx.json"

`

type Config struct {
	Core    *CoreConfig             `toml:"core"`
	Proxies map[string]*ProxyConfig `toml:"proxies"`
}

type CoreConfig struct {
	LogFile       string `toml:"log_file"`
	CloudSqlProxy string `toml:"cloud_sql_proxy"`
}

func NewConfig() *Config {
	return &Config{
		Core: &CoreConfig{
			LogFile:       "",
			CloudSqlProxy: BuiltinCloudSQLProxy,
		},
		Proxies: map[string]*ProxyConfig{},
	}
}

func (c *Config) Load(path string) error {
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

func (c *Config) SortedProxyKeys() []string {
	ret := []string{}
	for key, _ := range c.Proxies {
		ret = append(ret, key)
	}

	sort.Strings(ret)

	return ret
}

type ProxyConfig struct {
	Name    string `toml:"-"`
	Label   string `toml:"label"`
	Options string `toml:"options"`
}

func (c *ProxyConfig) NameForItem() string {
	if c.Label != "" {
		return c.Label
	}

	return c.Name
}

func (c *ProxyConfig) TooltipForItem() string {
	return c.Name
}
