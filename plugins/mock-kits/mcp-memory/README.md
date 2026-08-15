# MCP server: Memory (knowledge graph)

A mocked Memory MCP server: create entities and relations, add observations and search the graph. Gives an agent a deterministic memory contour to test recall logic against, with no shared state to reset between runs.

## Установка

```bash
mockarty-cli plugin pack examples/plugins/mcp-memory
mockarty-cli plugin install mockarty.mcp-memory-1.0.0.zip
mockarty-cli plugin enable mockarty.mcp-memory
```

Затем разверните кит в пространство имён — из каталога «Комплекты моков» или вызовом:

```bash
curl -X POST "$MOCKARTY_SERVER/api/v1/mock-kits/instantiate?namespace=$NS" \
  -H "Authorization: Bearer $MOCKARTY_TOKEN" -H 'Content-Type: application/json' \
  -d '{"kit":"mcp_memory"}'
```

## Инструменты

Агент подключается к `POST $MOCKARTY_SERVER/stubs/<пространство>` (JSON-RPC 2.0) и находит инструменты через `tools/list` — со схемами аргументов, включая `required`.

- `add_observations`
- `create_entities`
- `create_relations`
- `delete_entities`
- `read_graph`
- `search_nodes`

## Сценарии отказов

Вызовите инструмент с аргументом-триггером (описан в тексте ошибки), чтобы получить нужный отказ:

- `not-found`

## Лицензия

MIT
