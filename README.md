# Server Monitor

## Description

Server Monitor is a Go application designed to monitor the availability of server, port availability and service status of self and configured servers. It collects real-time data and store the metrics in the database. It notifies through telegram in case of anomalies.

## Development prerequisites

- Go (go1.23.4)
- PostgreSQL

## Local setup

Download the repo using,

```bash
git clone https://github.com/dharunvs/server-monitor.git
```

(or)

```bash
wget -O server-monitor.zip https://github.com/dharunvs/server-monitor/archive/refs/heads/main.zip
unzip server-monitor
```

#### Database

Create a database or use the `postgres` database.
Create the `monitoring_data` table use the query in `queries.sql`

#### Configuration

Rename or copy the `sample.config,json` to `config.json`
Modify the configuration accordingly

#### Run the application

Run the application using

```bash
go run cmd/main.go
```

build using

```bash
go build .
```

## Future Scope

- Database backup at regular intervals to configured backup servers
- Database replication checks at regular intervals to configured fallback servers
- Database restoration
- Email notification
