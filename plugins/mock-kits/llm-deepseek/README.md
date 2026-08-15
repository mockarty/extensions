# LLM provider: DeepSeek

A mocked DeepSeek API (OpenAI-compatible): chat with streaming, tool calling, a deepseek-reasoner reply carrying reasoning_content, balance, and the usual failure paths. Set the base URL to the mock and test without spending tokens.

## Установка

```bash
mockarty-cli plugin pack examples/plugins/llm-deepseek
mockarty-cli plugin install mockarty.llm-deepseek-1.0.0.zip
mockarty-cli plugin enable mockarty.llm-deepseek
```

Затем разверните кит в пространство имён — из каталога «Комплекты моков» или вызовом:

```bash
curl -X POST "$MOCKARTY_SERVER/api/v1/mock-kits/instantiate?namespace=$NS" \
  -H "Authorization: Bearer $MOCKARTY_TOKEN" -H 'Content-Type: application/json' \
  -d '{"kit":"llm_deepseek","prefix":"/api"}'
```

## Эндпоинты

- `GET /models`
- `GET /user/balance`
- `POST /chat/completions`

Отдаются по пути `$MOCKARTY_SERVER/stubs/<пространство>/<префикс><маршрут>`.
Стриминг включается самим запросом (`"stream": true`) — отдельный мок для этого не нужен.

## Сценарии отказов

Передайте имя модели `mock/<сценарий>`, чтобы получить нужный отказ:

- `content-filter`
- `context-length`
- `rate-limit`
- `reasoning`
- `server-error`
- `slow`
- `tool-calls`
- `truncated`
- `unauthorized`

## Лицензия

MIT
