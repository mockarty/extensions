# LLM provider: Qwen (DashScope)

A mocked Qwen API in DashScope's OpenAI-compatible mode: chat with streaming, tool calling, embeddings and the failure paths. Test a Qwen integration without a key or a quota.

## Установка

```bash
mockarty-cli plugin pack examples/plugins/llm-qwen
mockarty-cli plugin install mockarty.llm-qwen-1.0.0.zip
mockarty-cli plugin enable mockarty.llm-qwen
```

Затем разверните кит в пространство имён — из каталога «Комплекты моков» или вызовом:

```bash
curl -X POST "$MOCKARTY_SERVER/api/v1/mock-kits/instantiate?namespace=$NS" \
  -H "Authorization: Bearer $MOCKARTY_TOKEN" -H 'Content-Type: application/json' \
  -d '{"kit":"llm_qwen","prefix":"/api"}'
```

## Эндпоинты

- `GET /compatible-mode/v1/models`
- `POST /compatible-mode/v1/chat/completions`
- `POST /compatible-mode/v1/embeddings`

Отдаются по пути `$MOCKARTY_SERVER/stubs/<пространство>/<префикс><маршрут>`.
Стриминг включается самим запросом (`"stream": true`) — отдельный мок для этого не нужен.

## Сценарии отказов

Передайте имя модели `mock/<сценарий>`, чтобы получить нужный отказ:

- `content-filter`
- `context-length`
- `rate-limit`
- `server-error`
- `slow`
- `tool-calls`
- `truncated`
- `unauthorized`

## Лицензия

MIT
