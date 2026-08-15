# MCP server: Git

A mocked Git MCP server: inspect status and diffs, read the log, create branches and commit — including the errors an agent hits (not a repository, nothing staged, merge conflict). Test a code agent's git flow without touching a working tree.

## Установка

```bash
mockarty-cli plugin pack examples/plugins/mcp-git
mockarty-cli plugin install mockarty.mcp-git-1.0.0.zip
mockarty-cli plugin enable mockarty.mcp-git
```

Затем разверните кит в пространство имён — из каталога «Комплекты моков» или вызовом:

```bash
curl -X POST "$MOCKARTY_SERVER/api/v1/mock-kits/instantiate?namespace=$NS" \
  -H "Authorization: Bearer $MOCKARTY_TOKEN" -H 'Content-Type: application/json' \
  -d '{"kit":"mcp_git"}'
```

## Инструменты

Агент подключается к `POST $MOCKARTY_SERVER/stubs/<пространство>` (JSON-RPC 2.0) и находит инструменты через `tools/list` — со схемами аргументов, включая `required`.

- `git_commit`
- `git_create_branch`
- `git_diff_unstaged`
- `git_log`
- `git_status`

## Сценарии отказов

Вызовите инструмент с аргументом-триггером (описан в тексте ошибки), чтобы получить нужный отказ:

- `empty-message`
- `not-a-repo`

## Лицензия

MIT
