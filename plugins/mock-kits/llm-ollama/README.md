# LLM provider: Ollama (local)

A mocked Ollama server: the OpenAI-compatible chat endpoint with streaming and tool calling, plus the native /api/tags and /api/show. Lets CI exercise a local-model path without pulling multi-gigabyte weights.

## Установка

```bash
mockarty-cli plugin pack examples/plugins/llm-ollama
mockarty-cli plugin install mockarty.llm-ollama-1.0.0.zip
mockarty-cli plugin enable mockarty.llm-ollama
```

Затем разверните кит в пространство имён — из каталога «Комплекты моков» или вызовом:

```bash
curl -X POST "$MOCKARTY_SERVER/api/v1/mock-kits/instantiate?namespace=$NS" \
  -H "Authorization: Bearer $MOCKARTY_TOKEN" -H 'Content-Type: application/json' \
  -d '{"kit":"llm_ollama","prefix":"/api"}'
```

## Эндпоинты

- `GET /api/tags`
- `POST /api/show`
- `POST /v1/chat/completions`

Отдаются по пути `$MOCKARTY_SERVER/stubs/<пространство>/<префикс><маршрут>`.
Стриминг включается самим запросом (`"stream": true`) — отдельный мок для этого не нужен.

## Сценарии отказов

Передайте имя модели `mock/<сценарий>`, чтобы получить нужный отказ:

- `content-filter`
- `context-length`
- `down`
- `rate-limit`
- `server-error`
- `slow`
- `tool-calls`
- `truncated`
- `unauthorized`

## Лицензия

MIT
