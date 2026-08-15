# LLM provider: OpenAI

A complete OpenAI-compatible contour: chat completions (streaming and not), tool calling, embeddings, models — plus the failures that actually break integrations: 429 with Retry-After, 401, 500, context-length, truncation and content filtering. Point your client's base URL at the mock and test without a key or a bill.

## Установка

```bash
mockarty-cli plugin pack examples/plugins/llm-openai
mockarty-cli plugin install mockarty.llm-openai-1.0.0.zip
mockarty-cli plugin enable mockarty.llm-openai
```

Затем разверните кит в пространство имён — из каталога «Комплекты моков» или вызовом:

```bash
curl -X POST "$MOCKARTY_SERVER/api/v1/mock-kits/instantiate?namespace=$NS" \
  -H "Authorization: Bearer $MOCKARTY_TOKEN" -H 'Content-Type: application/json' \
  -d '{"kit":"openai_chat","prefix":"/api"}'
```

## Эндпоинты

- `GET /v1/models`
- `POST /v1/chat/completions`
- `POST /v1/embeddings`

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
