---
name: mockarty-plugin
description: Author, build, test and install a Mockarty plugin (extension) end-to-end. Use when the user wants to extend Mockarty — add a mock kit, a custom $.fake.* faker, a response transformer, a custom condition operator, a connector to an external tool, or a UI panel — and package it as an installable plugin. Covers scaffolding, the WASM guest contract, packing, and installing over CLI or the registry.
---

Mockarty plugins extend the product with **zero risk to the host**: declarative
content needs no code, and code runs only inside a locked WebAssembly sandbox (no
disk, no network, no host process). One manifest (`plugin.json`), one zip, install
from a file — air-gap friendly, nothing phones home.

Pick the mechanism that fits, copy its template, build if it has a `.wasm`, then
`pack` → `install` → `enable`. Everything below works fully offline; only
`install`/`enable` need a running Mockarty and an admin token.

## The six mechanisms

| Want to… | Mechanism | Code? |
|---|---|---|
| ship ready-made mocks (a whole API contour) | **mock kit** (`mock_kits`) | no — JSON |
| add a new `$.fake.<name>` generator | **WASM faker** (`wasm` @ `faker-provider`) | tiny WASM fn |
| rewrite a response body before it's sent | **WASM transformer** (`wasm` @ `response-transformer`) | tiny WASM fn |
| add a custom condition operator `plugin:<key>` | **WASM matcher** (`wasm` @ `condition-matcher`) | tiny WASM fn |
| pre-configure an external tool (Jira, GitLab, …) | **connector** (`connectors`) | no — JSON |
| seed another module (wiki/dashboards/issues/collections/TCM/contract) | **content pack** (`content_packs`: wiki/dashboard/issue/collection/tcm/contract) | no — JSON |
| relate entities to external systems (deep links) | **link type** (`links`) | no — JSON |
| add subscribable notification events | **event type** (`event_types`, `plugin.*` namespace) | no — JSON |
| add a dispatchable runner job kind | **task type** (`task_types`, `plugin-*`, pair with an external runner) | no — JSON |
| add a page/panel/tab/command in the UI | **UI panel** (`ui`: sidebar / page-slot / entity-tab / command / settings-section) | HTML |

## Fastest path (CLI scaffolds everything)

```bash
mockarty-cli plugin create my-plugin --template mock-kit     # or wasm-faker / connector / ui-panel
cd my-plugin
# edit plugin.json (and guest/main.go if it's a WASM template)
./build.sh                      # only for WASM templates — needs TinyGo (tinygo.org)
mockarty-cli plugin pack .      # → my-plugin-1.0.0.zip  (prints sha256)
export MOCKARTY_SERVER=http://localhost:5770
export MOCKARTY_TOKEN=mk_...    # Admin → API tokens
mockarty-cli plugin install my-plugin-1.0.0.zip
mockarty-cli plugin enable <plugin-id>
```

`mockarty-cli plugin dev .` hot-reloads on every file change while you iterate.

**Already have mocks elsewhere?** Convert a WireMock / Postman / Mockoon library
into a plugin in one command:

```bash
mockarty-cli plugin import ./mappings --from wiremock --id acme.legacy   # or postman / mockoon
```

## plugin.json skeleton

```json
{
  "id": "acme.my-plugin",
  "name": "My plugin",
  "version": "1.0.0",
  "description": "One clear sentence — shown in the catalogue.",
  "author": {"name": "You", "url": "https://example.com"},
  "license": "MIT",
  "contributes": { "...": "one or more of the blocks below" }
}
```

`id` is `vendor.name`, globally unique. Bump `version` (semver) on every change.

## Templates by mechanism

### 1. Mock kit (declarative — no code)
```json
"contributes": { "mock_kits": [
  { "key": "users_api", "name": "Users API", "description": "Starter contour.",
    "mocks": [
      {"route": "/users", "method": "GET", "status_code": 200, "description": "List",
       "body": {"users": [{"id": "$.fake.UUID", "name": "$.fake.Name", "email": "$.fake.Email"}]}},
      {"route": "/users", "method": "POST", "status_code": 201, "body": {"id": "$.fake.UUID", "created": true}}
    ] } ] }
```
Bodies support the same `$.fake.*` helpers a normal mock does (UUID, Name, Email,
Company, City, Price, IPv4, `IntBetween(min,max)`, …). After install+enable,
instantiate the kit into a namespace from the Mock Kits catalogue.

### 2. WASM faker — a new `$.fake.<name>`
```json
"contributes": { "wasm": [
  { "point": "faker-provider", "module": "faker.wasm", "fn": "faker", "exports": ["my_id"] } ] }
```
Guest ABI (TinyGo, `-target=wasm-unknown`, build tag `tinygo`): export `memory`,
`mk_alloc(size)->ptr`, and `faker(ptr,len)->i64` packing `(outPtr<<32|outLen)`.
Input `{"name":"my_id"}` → output `{"value":"..."}`. **Zero-alloc** (freestanding
GC leaks): fixed buffers, no `make`/`+`/string concat. See
`examples/plugins/ru-fakers` for the full pattern.

### 3. WASM response transformer — rewrite a body
```json
"contributes": { "wasm": [
  { "point": "response-transformer", "module": "transform.wasm", "fn": "transform", "exports": ["envelope"] } ] }
```
`transform(ptr,len)->i64` gets the raw body bytes, returns the new body. Opt in
per mock with the tag `transform:envelope`. Host caps input at 1 MiB; a broken
transformer is skipped (original body sent). Example: `examples/plugins/response-transformer`.

### 4. WASM condition matcher — a custom operator `plugin:<key>`
```json
"contributes": { "wasm": [
  { "point": "condition-matcher", "module": "match.wasm", "fn": "match", "exports": ["luhn"] } ] }
```
`match(ptr,len)->i64` gets `{"key","actual","expected"}`, returns `{"match":bool}`.
Use it as a mock condition: `assertAction: "plugin:luhn"` (path picks the value,
value is the operand). A missing/disabled matcher never matches (safe default).
Example: `examples/plugins/condition-matcher`.

### 5. Connector — pre-configure an external tool (no code)
```json
"contributes": { "connectors": [
  { "key": "jira_cloud", "name": "Jira Cloud", "kind": "jira", "icon": "bug-ant",
    "config_template": {"base_url": "https://your-org.atlassian.net"} } ] }
```
`kind` must be a built-in adapter: `jira`, `github`, `gitlab`, `jenkins`, `linear`,
`allure`, `testrail`, `webhook_generic`. **Never put secrets in the manifest** —
credentials are added per-namespace when the integration is created. Example:
`examples/plugins/jira-connector`.

### 6. UI panel — a page in the sidebar
```json
"contributes": { "ui": [
  { "point": "sidebar", "id": "ops", "title": "Ops", "icon": "chart-bar", "panel": "panel.html" } ] }
```
`panel.html` (bundled) renders in a sandboxed iframe — strict CSP, no host access;
talk to the host via `postMessage`. Points: `sidebar`, `page-slot` (needs
`target`), `entity-tab` (target = `issue` / `mock` / `case`), `settings-section`. Example: `examples/plugins/ops-panel`.

## Multi-tenancy (how installs work)

An admin installs a plugin once for the whole app; each namespace can then
**enable/disable** it and **override its settings** independently. A plugin can
ship `settings_schema` (JSON Schema) → the host renders a config form. Default a
plugin off with `"default_off": true` for opt-in features.

## Checklist before you ship

- [ ] `mockarty-cli plugin inspect my-plugin-1.0.0.zip` prints the contributions you expect.
- [ ] WASM: `.wasm` built and its exports listed match what the guest exports.
- [ ] Install → enable → the mechanism works (kit instantiates / faker resolves / etc.).
- [ ] No secrets in `plugin.json`.
- [ ] Bumped `version`.

Full authoring guide (every field, both audiences): **/docs/plugins-authoring**.
