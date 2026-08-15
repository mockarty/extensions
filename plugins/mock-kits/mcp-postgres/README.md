# MCP server: PostgreSQL

A mocked PostgreSQL MCP server: explore schemas and tables, describe columns and run read-only queries — including the refusals a safe server enforces (write statements rejected) and real SQL errors. Test a data agent without a database.

## Установка

```bash
mockarty-cli plugin pack examples/plugins/mcp-postgres
mockarty-cli plugin install mockarty.mcp-postgres-1.0.0.zip
mockarty-cli plugin enable mockarty.mcp-postgres
```

Затем разверните кит в пространство имён — из каталога «Комплекты моков» или вызовом:

```bash
curl -X POST "$MOCKARTY_SERVER/api/v1/mock-kits/instantiate?namespace=$NS" \
  -H "Authorization: Bearer $MOCKARTY_TOKEN" -H 'Content-Type: application/json' \
  -d '{"kit":"mcp_postgres"}'
```

## Инструменты

Агент подключается к `POST $MOCKARTY_SERVER/stubs/<пространство>` (JSON-RPC 2.0) и находит инструменты через `tools/list` — со схемами аргументов, включая `required`.

- `execute_sql`
- `explain_query`
- `get_object_details`
- `list_objects`
- `list_schemas`

## Сценарии отказов

Вызовите инструмент с аргументом-триггером (описан в тексте ошибки), чтобы получить нужный отказ:

- `not-found`
- `read-only`
- `sql-error`

## Лицензия

MIT
