# tcpmon

A simple TCP monitor for Linux.

## Installation

```bash
# build the RPM
make package

# install
rpm -Uvh tcpmon-<version>.el7.x86_64.rpm
```

## Getting started

Start collecting monitoring data:

```bash
tcpmon start
```

Collection metrics snapshot:

```bash
curl -JfSsLO http://127.0.0.1:6789/backup
```

Export metrics in line protocol:

```bash
tcpmon export -o metrics.txt <backup-dir>
```

Then import `metrics.txt` in InfluxDB and find out what is going wrong.

## Configuration

Config file located at `$HOME/.tcpmon/config.yaml` (Development) or `/etc/tcpmon/config.yaml` (Production)

## Development

```bash
# build a binary
make

# build binary for Linux
make build-linux
```

## License

MIT
