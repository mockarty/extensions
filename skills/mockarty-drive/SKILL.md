---
name: mockarty-drive
description: Drive Mockarty end-to-end as an autonomous backend tester — create and manage API mocks (HTTP/gRPC/GraphQL/SOAP/SSE/WebSocket), run functional collections, load/perf tests, fuzzing, chaos and contract checks, manage test cases, and read reports. Use when the user wants to mock a backend, exercise an API, or run any kind of testing through Mockarty over MCP (preferred) or REST.
---

Mockarty is a backend-testing platform: API mocking + functional/load/fuzz/chaos/
contract testing + test-case management, all driveable by an agent. Prefer the
**MCP tools** (`mockarty` server) — they carry namespace defaults, licence
filtering and safety annotations. Fall back to REST only when MCP isn't wired.

## Conventions (read first)

- **Namespace**: omit it and your pinned/default namespace is used. Pass one only
  to target another namespace you can access.
- **Names → ids**: never guess ids. Resolve a human name with
  `search_entities {type, query}` (type is required), then call the matching
  `get_/run_/update_` tool.
- **Long-running ops** (`run_perf_test`, `start_fuzzing`, `run_collection`,
  `run_test_plan`, `chaos_run_experiment`, `start_recording`) return an id
  immediately — poll `wait_for_task {task_type, task_id}` or the matching `get_`
  tool until a terminal status. Never busy-loop.
- **Fakers** take no args (`$.fake.UUID`, `$.fake.Email`, `$.fake.PositiveInt`),
  except `$.fake.IntBetween(min,max)`.
- **Errors** come back as structured JSON — read `offendingField`, `hint`,
  `example` before retrying. Never fabricate ids, statuses or payloads.

## Mock a backend

```
create_mock { route:"/users/:id", method:"GET", responseHttp:{
  code:200, body:{ id:"$.fake.UUID", name:"$.fake.Name", email:"$.fake.Email" } } }
test_mock { ... }          # exercise it without leaving Mockarty
list_mocks {}              # see what exists; search_entities to resolve one by name
```

Bootstrap a whole contour from a spec instead of hand-writing mocks:
`generate_from_openapi`, `import_openapi`, `import_postman`, `import_har`,
`import_wsdl`, `import_graphql`, `mock_bootstrap_from_spec`. To seed ready-made
content from an installed plugin: `mock_kit_list` → `mock_kit_instantiate`,
`content_pack_list` → `content_pack_instantiate`.

Dynamic responses use `$.fake.*` (data), `$.jsonPath(...)` (echo request), stores
(Global/Chain/Mock) and condition matching — see the mock a real request hits.

## Run the tests

| Kind | Start | Read result |
|------|-------|-------------|
| Functional collection | `run_collection` | `get_test_run_report` |
| Load / perf | `run_perf_test` (or `create_load_campaign`) | `get_perf_result` / `get_load_campaign_report` |
| Fuzzing | `start_fuzzing` (`import_fuzz_from_openapi` to seed) | `list_fuzz_findings`, `export_fuzz_findings` |
| Chaos | `chaos_run_experiment` / `chaos_run_preset` | `chaos_get_experiment_report` |
| Contract | `contract_validate_mocks`, `contract_verify_provider`, `contract_pact_can_i_deploy` | inline |
| Test plan (mixed suite) | `run_test_plan` | `get_test_plan_run_report` / `_unified` / `_junit` |

Always `wait_for_task` (or poll the `get_` tool) to a terminal status before
reading the report. Reports come in JSON, Markdown, HTML and JUnit — pick JUnit
for CI, Markdown to summarise for a human.

## Manage test cases (TCM)

```
tcm_cases_create { name, priority, steps:[{action, expectedResult}] }
tcm_cases_list / tcm_cases_get / tcm_cases_run
create_test_plan → add cases → run_test_plan
```
Link a verification back to its issue in the tracker with
`issuetracker_issue_link_test` so traceability holds.

## Record → replay

`start_recording` (proxy real traffic) → `stop_recording` →
`create_mocks_from_recording` / `export_recording_collection` /
`export_recording_perf_script`. Turns live traffic into mocks + a load script in
two calls.

## A good autonomous loop

1. `current_namespace` — know where you are.
2. Create/import the mocks the system-under-test needs.
3. Run the relevant suite(s); `wait_for_task` to completion.
4. Read the report; if there are failures/findings, triage
   (`triage_fuzz_finding`, `security_triage_finding`) and fix the mock/spec.
5. Re-run until green; record the outcome as a TCM case linked to the work item.

## REST fallback

Every tool maps to `POST/GET /api/v1/...` on the admin host with an
`Authorization: Bearer mk_...` token (Admin → API tokens). Example:

```bash
curl -sS -X POST "$MOCKARTY/api/v1/mocks?namespace=$NS" \
  -H "Authorization: Bearer $TOKEN" -H 'Content-Type: application/json' \
  -d '{"name":"users","route":"/users","method":"GET","responseType":"http",
       "responseHttp":{"code":200,"body":"{\"ok\":true}"}}'
```

`http_request` (MCP) is for hitting mock endpoints or external URLs to test them —
NOT the admin REST API (loopback/private targets are blocked). Use the dedicated
tools for admin operations.

Full product docs (the agent's source of truth, also RAG-indexed): the **/docs**
site on your Mockarty host.
