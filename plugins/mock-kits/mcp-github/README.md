# MCP server: GitHub

A mocked GitHub MCP server: search code, read files, open and comment on issues, list and create pull requests — with the failures an agent must handle (repository not found, insufficient permissions, rate limit). Drive an agent through a whole GitHub workflow without a token or a real repository.

## Установка

```bash
mockarty-cli plugin pack examples/plugins/mcp-github
mockarty-cli plugin install mockarty.mcp-github-1.0.0.zip
mockarty-cli plugin enable mockarty.mcp-github
```

Затем разверните кит в пространство имён — из каталога «Комплекты моков» или вызовом:

```bash
curl -X POST "$MOCKARTY_SERVER/api/v1/mock-kits/instantiate?namespace=$NS" \
  -H "Authorization: Bearer $MOCKARTY_TOKEN" -H 'Content-Type: application/json' \
  -d '{"kit":"mcp_github"}'
```

## Инструменты

Агент подключается к `POST $MOCKARTY_SERVER/stubs/<пространство>` (JSON-RPC 2.0) и находит инструменты через `tools/list` — со схемами аргументов, включая `required`.

- `add_issue_comment`
- `create_issue`
- `create_pull_request`
- `get_file_contents`
- `list_pull_requests`
- `search_repositories`

## Сценарии отказов

Вызовите инструмент с аргументом-триггером (описан в тексте ошибки), чтобы получить нужный отказ:

- `conflict`
- `forbidden`
- `not-found`

## Лицензия

MIT
