# MCP server: SQLite

A mocked SQLite MCP server: list and describe tables, read and write rows, and record business insights — with genuine SQLite error text for the failure paths. Useful for testing an agent against a small embedded database it cannot corrupt.

## Установка

```bash
mockarty-cli plugin pack examples/plugins/mcp-sqlite
mockarty-cli plugin install mockarty.mcp-sqlite-1.0.0.zip
mockarty-cli plugin enable mockarty.mcp-sqlite
```

Затем разверните кит в пространство имён — из каталога «Комплекты моков» или вызовом:

```bash
curl -X POST "$MOCKARTY_SERVER/api/v1/mock-kits/instantiate?namespace=$NS" \
  -H "Authorization: Bearer $MOCKARTY_TOKEN" -H 'Content-Type: application/json' \
  -d '{"kit":"mcp_sqlite"}'
```

## Инструменты

Агент подключается к `POST $MOCKARTY_SERVER/stubs/<пространство>` (JSON-RPC 2.0) и находит инструменты через `tools/list` — со схемами аргументов, включая `required`.

- `append_insight`
- `create_table`
- `describe_table`
- `list_tables`
- `read_query`
- `write_query`

## Сценарии отказов

Вызовите инструмент с аргументом-триггером (описан в тексте ошибки), чтобы получить нужный отказ:

- `not-found`
- `wrong-tool`

## Лицензия

MIT
