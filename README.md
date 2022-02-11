# Alertmanager-telegram
Alertmanager telegram receiver

[CLI docs](CLI.md)

## Example alerts:
![Alerts example](alerts.png)

## Docker images
 * amd64 `ghcr.io/a1fred/alertmanager-telegram:[latest|<version>]`
 * arm64/v8 `ghcr.io/a1fred/alertmanager-telegram-arm64v8:[latest|<version>]`

Get help
```shell
$ docker run --rm ghcr.io/a1fred/alertmanager-telegram --help
alertmanager-telegram@0.0.1-3a73ee0ea1c94ab1fc2ea03ab56774300d7c9d2a-2022-02-10T22:33:22Z
Usage:
  alertmanager-telegram [OPTIONS] <daemon>

Help Options:
  -h, --help  Show this help message

Available commands:
  daemon  Daemon
```

Example run
```shell
$ docker run --rm \
    -e TELEGRAM_TOKEN="<bot token>" \
    -e TELEGRAM_CHAT_ID="<your chat id>" \
    ghcr.io/a1fred/alertmanager-telegram
```

## Healthcheck
You can use `/health` endpoint, always responds `200`, can be useful for health-checking.

## Metrics
prometheus metrics available on `/metrics` endpoint
```
# HELP alertmanager_telegram_alerts_received Number of alerts received
# TYPE alertmanager_telegram_alerts_received counter
alertmanager_telegram_alerts_received 3
# HELP alertmanager_telegram_messages_sending_error Number of errors message sending to telegram recipients
# TYPE alertmanager_telegram_messages_sending_error counter
alertmanager_telegram_messages_sending_error 0
# HELP alertmanager_telegram_messages_sent Number of messages sent to telegram recipients
# TYPE alertmanager_telegram_messages_sent counter
alertmanager_telegram_messages_sent 3
```

# Development
 * `make build` - build local binary
 * `make alertmanager_run` - run local instance of alertmanager
 * `make alert` - send dummy alert to alertmanager
 * run `build/alertmanager-telegram --telegram.token=<test bot token> -r=<your telegram id>` for receive alerts
