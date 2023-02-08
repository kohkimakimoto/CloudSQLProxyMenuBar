# CloudSQLProxyMenuBar

CloudSQLProxyMenuBar displays a menu to manage [Cloud SQL Auth Proxy](https://cloud.google.com/sql/docs/mysql/sql-proxy) processes on macOS menu bar.

![screenshot.png](https://raw.githubusercontent.com/kohkimakimoto/CloudSQLProxyMenuBar/master/screenshot.png)

## Installation

[Download the latest version](https://github.com/kohkimakimoto/CloudSQLProxyMenuBar/releases/latest)

## Configuration

CloudSQLProxyMenuBar loads configuration from `$HOME/.cloudsqlproxymenubar/config.toml`. See the following example:

```toml
#
# core is the section of CloudSQLProxyMenuBar global config.
#
[core]
# Required: The path to 'cloud_sql_proxy' command.
# If you are not familiar with cloud_sql_proxy, please read the document: https://cloud.google.com/sql/docs/mysql/sql-proxy
cloud_sql_proxy = "/path/to/cloud_sql_proxy"

# Optional: The log file path.
# The default is '$HOME/.cloudsqlproxymenubar/output.log'
log_file = "/path/to/logfile"

#
# proxies.xxx is the section of the Cloud SQL Proxy settings.
#
[proxies.cloudsqlinstance1]
# Optional: The text is displayed on the menu item.
# The default is the same as 'XXX' part of 'proxies.XXX'.
label = "proxy-to-cloudsqlinstance1"

# Required: The command line arguments passed to 'cloud_sql_proxy' command.
arguments = "-dir=/cloudsql -instances=yourproject:asia-northeast1:cloudsqlinstance1 -credential_file=/path/to/service_account.json"

# You can set proxy config multiple times.
# [proxies.cloudsqlinstance2]
# ...
```

## Author

Kohki Makimoto <kohki.makimoto@gmail.com>

## License

The MIT License (MIT)
