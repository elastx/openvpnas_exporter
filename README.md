# OpenVPN AS metrics exporter for Prometheus

A Prometheus exporter that makes calls to the XML-RPC unix domain socket exposed by the [OpenVPN AS](https://openvpn.net/) service, and generates metrics from the responses.

Current metrics include:

* Server status
  * Server version
  * Number of connected clients and max allowed concurrect connections
* Subscription status
  * Number of current and maximum concurrent connections
  * Subscription last successful sync and next update sync timestamp

Based on [openvpnas-exporter](https://github.com/rossigee/openvpnas-exporter).


## Exposed metrics example

```
# HELP openvpnas_server_connected_clients Number of currently connected clients to the server.
# TYPE openvpnas_server_connected_clients gauge
openvpnas_server_connected_clients 1
# HELP openvpnas_server_connected_clients_limit Server concurrent client connection limit.
# TYPE openvpnas_server_connected_clients_limit gauge
openvpnas_server_connected_clients_limit 10
# HELP openvpnas_server_version_info Contains OpenVPN AS server version info
# TYPE openvpnas_server_version_info gauge
openvpnas_server_version_info{version="2.12.3 (build 76774795)"} 1
# HELP openvpnas_subscription_connected_clients Total number of client connections currently in use across the OpenVPN AS subscription.
# TYPE openvpnas_subscription_connected_clients gauge
openvpnas_subscription_connected_clients 3
# HELP openvpnas_subscription_connected_clients_limit Maximum number of concurrent client connections allowed by the OpenVPN AS subscription.
# TYPE openvpnas_subscription_connected_clients_limit gauge
openvpnas_subscription_connected_clients_limit 10
# HELP openvpnas_subscription_status_last_update_time_seconds UNIX timestamp when the OpenVPN AS subscription was last synced.
# TYPE openvpnas_subscription_status_last_update_time_seconds gauge
openvpnas_subscription_status_last_update_time_seconds 1.708366094e+09
# HELP openvpnas_subscription_status_next_update_time_seconds UNIX timestamp of the next planned OpenVPN AS subscription sync.
# TYPE openvpnas_subscription_status_next_update_time_seconds gauge
openvpnas_subscription_status_next_update_time_seconds 1.708366294e+09
# HELP openvpnas_up Whether scaping OpenVPN AS metrics was successful.
# TYPE openvpnas_up gauge
openvpnas_up 1
```


## Usage

```sh
usage: openvpnas_exporter [<flags>]


Flags:
  --[no-]help                   Show context-sensitive help (also try --help-long and --help-man).
  --web.listen-address=":9176"  Address to listen on for web interface and telemetry.
  --web.telemetry-path="/metrics"  
                                Path under which to expose metrics.
  --openvpnas.xmlrpc-path="/usr/local/openvpn_as/etc/sock/sagent.localroot"  
                                Path to the XML-RPC unix domain socket file.
  --[no-]version                Show application version.
```

By default the socket file is owned by root so you probably need to run the exporter as root (for now).


## Get a standalone executable binary

You can download the pre-compiled binaries from the [releases page](https://github.com/elastx/openvpnas_exporter/releases).


## Contributions/development

Please run `gofmt` on your code to fix formatting/styling, ie. `gofmt -w .` to recursively fix all formatting in the repository.


### Dependencies

Go is fairly smart about this, if you add a import to a program and run `go build path/to/main.go` it will automatically update `go.mod`.

Still, it is a good idea to run `go mod tidy` before committing your changes, this will ensure `go.mod` and `go.sum` are in order.


### Commit messages

If you want your changes to show up properly in the release notes, prefix your commit messages with:

- build
- doc
- feat
- fix
- sec

Eg.

> feat: add openvpnas_up gauge

See `.changelog.groups` in `.goreleaser` for details / full regex.


### Releases

To make a new release you need to create a new tag in the semver format, ie tag name will look something like `v0.0.1`.

```shell
git checkout master
git tag -a v0.0.1 -m "My message (not shown anywhere)"
git push origin v0.0.1
```
