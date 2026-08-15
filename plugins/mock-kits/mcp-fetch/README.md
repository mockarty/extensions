# MCP server: Fetch (web)

A mocked Fetch MCP server: retrieve a web page as markdown, with the failures an agent has to survive — 404, timeout, and a robots.txt refusal. Lets you test retrieval logic offline and deterministically.

## Установка

```bash
mockarty-cli plugin pack examples/plugins/mcp-fetch
mockarty-cli plugin install mockarty.mcp-fetch-1.0.0.zip
mockarty-cli plugin enable mockarty.mcp-fetch
```

Затем разверните кит в пространство имён — из каталога «Комплекты моков» или вызовом:

```bash
curl -X POST "$MOCKARTY_SERVER/api/v1/mock-kits/instantiate?namespace=$NS" \
  -H "Authorization: Bearer $MOCKARTY_TOKEN" -H 'Content-Type: application/json' \
  -d '{"kit":"mcp_fetch"}'
```

## Инструменты

Агент подключается к `POST $MOCKARTY_SERVER/stubs/<пространство>` (JSON-RPC 2.0) и находит инструменты через `tools/list` — со схемами аргументов, включая `required`.

- `fetch`

## Сценарии отказов

Вызовите инструмент с аргументом-триггером (описан в тексте ошибки), чтобы получить нужный отказ:

- `not-found`
- `robots`
- `timeout`

## Лицензия

MIT
