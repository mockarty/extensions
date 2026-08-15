# MCP server: Playwright (browser)

A mocked Playwright MCP server: navigate, snapshot, click, type, screenshot, evaluate and close — with real input schemas so an agent can drive it from tools/list alone. Includes failure cases (element not found, navigation timeout) so you can test how your agent recovers without launching a browser.

## Установка

```bash
mockarty-cli plugin pack examples/plugins/mcp-playwright
mockarty-cli plugin install mockarty.mcp-playwright-1.0.0.zip
mockarty-cli plugin enable mockarty.mcp-playwright
```

Затем разверните кит в пространство имён — из каталога «Комплекты моков» или вызовом:

```bash
curl -X POST "$MOCKARTY_SERVER/api/v1/mock-kits/instantiate?namespace=$NS" \
  -H "Authorization: Bearer $MOCKARTY_TOKEN" -H 'Content-Type: application/json' \
  -d '{"kit":"mcp_playwright"}'
```

## Инструменты

Агент подключается к `POST $MOCKARTY_SERVER/stubs/<пространство>` (JSON-RPC 2.0) и находит инструменты через `tools/list` — со схемами аргументов, включая `required`.

- `browser_click`
- `browser_close`
- `browser_evaluate`
- `browser_navigate`
- `browser_snapshot`
- `browser_take_screenshot`
- `browser_type`

## Сценарии отказов

Вызовите инструмент с аргументом-триггером (описан в тексте ошибки), чтобы получить нужный отказ:

- `not-found`
- `timeout`

## Лицензия

MIT
