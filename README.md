# tcpmon

A simple TCP monitor for Linux.

## Installation

```bash
make rpm
```

## Usage

Start collecting monitoring data:

```bash
tcpmon start
```

Collection metrics snapshot:

```bash
curl -JfSsLO http://<ip>:6789/backup
```

Export metrics:

```bash
tcpmon export --format tsdb <backup-file> > db.txt
```

Then import `db.txt` in OpenTSDB, add it as Grafana DataSource, find out what is wrong.

## Configuration

Config file located at `$HOME/.tcpmon/config.yaml` (Development) or `/etc/tcpmon/config.yaml` (Production)

## Development

```bash
make
```

## License

MIT
