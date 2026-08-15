# MCP server: Brave Search

A mocked Brave Search MCP server: web and local search with realistic result shapes, plus the quota error and the empty-result case an agent must handle gracefully. Test search-and-summarise flows deterministically, with no API key and no per-query cost.

## Установка

```bash
mockarty-cli plugin pack examples/plugins/mcp-brave-search
mockarty-cli plugin install mockarty.mcp-brave-search-1.0.0.zip
mockarty-cli plugin enable mockarty.mcp-brave-search
```

Затем разверните кит в пространство имён — из каталога «Комплекты моков» или вызовом:

```bash
curl -X POST "$MOCKARTY_SERVER/api/v1/mock-kits/instantiate?namespace=$NS" \
  -H "Authorization: Bearer $MOCKARTY_TOKEN" -H 'Content-Type: application/json' \
  -d '{"kit":"mcp_brave_search"}'
```

## Инструменты

Агент подключается к `POST $MOCKARTY_SERVER/stubs/<пространство>` (JSON-RPC 2.0) и находит инструменты через `tools/list` — со схемами аргументов, включая `required`.

- `brave_local_search`
- `brave_web_search`

## Сценарии отказов

Вызовите инструмент с аргументом-триггером (описан в тексте ошибки), чтобы получить нужный отказ:

- `no-results`
- `rate-limit`

## Лицензия

MIT
