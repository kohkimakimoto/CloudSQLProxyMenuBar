# CloudSQLProxyMenuBar

CloudSQLProxyMenuBar displays a menu to manage [Cloud SQL Proxy](https://cloud.google.com/sql/docs/mysql/sql-proxy) processes on macOS menu bar.

![screenshot.png](https://raw.githubusercontent.com/kohkimakimoto/CloudSQLProxyMenuBar/master/screenshot.png)

## Installation

[Download latest version](https://github.com/kohkimakimoto/CloudSQLProxyMenuBar/releases/latest)

## Configuration

CloudSQLProxyMenuBar loads configuration from `$HOME/.cloudsqlproxymenubar/config.toml`. See the following example:

```toml
#
# core is the section of CloudSQLProxyMenuBar global config.
#
[core]
# optional: The path of `cloud_sql_proxy` command.
# If you do not set it, CloudSQLProxyMenuBar uses builtin `cloud_sql_proxy` command.
cloud_sql_proxy = "/usr/local/bin/cloud_sql_proxy"

# optional: The log file path.
# The default is `$HOME/.cloudsqlproxymenubar/output.log`
log_file = "/path/to/logfile"

#
# proxies.xxx is the section of the Cloud SQL Proxy settings.
#
[proxies.yourinstance1]
# optional: The text is displayed on the menu item.
# The default is the same as `XXX` of `proxies.XXX`.
label = "proxy-to-yourinstance1"
# required: The command line options of `cloud_sql_proxy` command.
options = "-dir=/cloudsql -instances=yourproject:asia-northeast1:yourinstance1"

# You can set proxy config multiple times.
# [proxies.yourinstance2]
# ...
```

## Author

Kohki Makimoto <kohki.makimoto@gmail.com>

## License

The MIT License (MIT)
