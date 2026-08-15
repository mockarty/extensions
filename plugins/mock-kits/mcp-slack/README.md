# MCP server: Slack

A mocked Slack MCP server: list channels, post and reply to messages, add reactions and look up users — plus the errors a bot really hits (not in channel, channel archived, rate limited). Test a Slack agent without a workspace.

## Установка

```bash
mockarty-cli plugin pack examples/plugins/mcp-slack
mockarty-cli plugin install mockarty.mcp-slack-1.0.0.zip
mockarty-cli plugin enable mockarty.mcp-slack
```

Затем разверните кит в пространство имён — из каталога «Комплекты моков» или вызовом:

```bash
curl -X POST "$MOCKARTY_SERVER/api/v1/mock-kits/instantiate?namespace=$NS" \
  -H "Authorization: Bearer $MOCKARTY_TOKEN" -H 'Content-Type: application/json' \
  -d '{"kit":"mcp_slack"}'
```

## Инструменты

Агент подключается к `POST $MOCKARTY_SERVER/stubs/<пространство>` (JSON-RPC 2.0) и находит инструменты через `tools/list` — со схемами аргументов, включая `required`.

- `slack_add_reaction`
- `slack_get_channel_history`
- `slack_get_users`
- `slack_list_channels`
- `slack_post_message`
- `slack_reply_to_thread`

## Сценарии отказов

Вызовите инструмент с аргументом-триггером (описан в тексте ошибки), чтобы получить нужный отказ:

- `archived`
- `invalid-emoji`
- `not-in-channel`

## Лицензия

MIT
