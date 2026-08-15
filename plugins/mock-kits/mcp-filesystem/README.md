# MCP server: Filesystem

A mocked Filesystem MCP server: read, write, edit, list and search files, with the sandbox refusals a real one enforces (path outside the allowed directories, file not found). Test an agent's file-handling logic without letting it near a real disk.

## Установка

```bash
mockarty-cli plugin pack examples/plugins/mcp-filesystem
mockarty-cli plugin install mockarty.mcp-filesystem-1.0.0.zip
mockarty-cli plugin enable mockarty.mcp-filesystem
```

Затем разверните кит в пространство имён — из каталога «Комплекты моков» или вызовом:

```bash
curl -X POST "$MOCKARTY_SERVER/api/v1/mock-kits/instantiate?namespace=$NS" \
  -H "Authorization: Bearer $MOCKARTY_TOKEN" -H 'Content-Type: application/json' \
  -d '{"kit":"mcp_filesystem"}'
```

## Инструменты

Агент подключается к `POST $MOCKARTY_SERVER/stubs/<пространство>` (JSON-RPC 2.0) и находит инструменты через `tools/list` — со схемами аргументов, включая `required`.

- `edit_file`
- `get_file_info`
- `list_allowed_directories`
- `list_directory`
- `read_text_file`
- `search_files`
- `write_file`

## Сценарии отказов

Вызовите инструмент с аргументом-триггером (описан в тексте ошибки), чтобы получить нужный отказ:

- `denied`
- `not-found`

## Лицензия

MIT
