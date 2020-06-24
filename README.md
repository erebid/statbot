# Statbot

A simple Discord bot that logs messages to a Postgres database.

I personally use it for having a Grafana dashboard with server stats.

```
go get github.com/erebid/statbot/cmd/statbot
statbot -t "discord bot token" -d "postgres database to connect to"
```
The option passed into -d can either be a URL (like `postgresql://statbot@/statbot`) or a DSN connection string (like `host=/var/run/postgresql user=statbot database=statbot`).

Statbot is free software, and its license is included in [LICENSE](./LICENSE).