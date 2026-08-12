# Mockarty Extensions

The official extension catalogue for [Mockarty](https://mockarty.ru) — ready-made
**plugins** (mock kits, custom fakers, response transformers, condition
matchers, connectors, content packs, UI panels, entity links, event & task
types) and **agent skills** for AI coding assistants.

Everything here installs **offline** into any Mockarty: a plugin is a small
`.zip` you upload; nothing phones home. Code-carrying plugins run in a locked
WebAssembly sandbox; UI panels render in a sandboxed iframe. See the product
docs: *Plugins*, *Writing plugins* and *Plugin operations*.

## Use this repo as your marketplace

Point your Mockarty at this repository's index and the whole catalogue appears
under **Admin → Plugins**, installable in one click with sha256 verification:

```bash
MOCKARTY_PLUGINS_REGISTRY_URL=https://raw.githubusercontent.com/mockarty/extensions/main/index.json
```

Air-gapped? Download a release zip from [Releases](../../releases) and install
it as a file — same result.

## Install a plugin

| How | Command / path |
|-----|----------------|
| **UI** | Admin → Plugins → *Install plugin…* (file) or *Install from URL…* (paste a release link, optional sha256 pin) |
| **CLI** | `mockarty-cli plugin install <file.zip \| https://…/file.zip \| plugin-id>` then `mockarty-cli plugin enable <id>` |
| **REST** | `POST /api/v1/plugins` (zip body) · `POST /api/v1/plugins/install-url` · `POST /api/v1/plugin-registry/install` |
| **AI agent (MCP)** | `plugin_registry_list` → `plugin_registry_install` → `plugin_enable` (or `plugin_install_url`) |

A fresh install lands **disabled** — enable it globally, then per namespace if
you use multi-tenancy.

## Catalogue

| Folder | What's inside |
|--------|---------------|
| [`plugins/mock-kits/`](plugins/mock-kits/) | Ready-made mock contours: whole API domains (users, e-commerce, payments, analytics, …) you instantiate into a namespace |
| [`plugins/fakers/`](plugins/fakers/) | Custom `$.fake.*` generators with **valid check digits** — RU (ИНН/СНИЛС/ОГРН), EU (IBAN/VAT), US (SSN/ABA), finance, crypto, network, datetime, identifiers, barcodes (EAN/UPC/ISBN) |
| [`plugins/transformers/`](plugins/transformers/) | WASM response transformers — rewrite a mock's body inline (envelopes, masking) |
| [`plugins/matchers/`](plugins/matchers/) | WASM condition matchers — custom `plugin:<key>` operators (Luhn, …) |
| [`plugins/connectors/`](plugins/connectors/) | Pre-configured bindings to external tools: Jira, GitHub, GitLab, Jenkins, Linear, Allure, TestRail, generic webhook |
| [`plugins/content-packs/`](plugins/content-packs/) | Seed other modules: team wiki (runbooks/ADR/PRD), Jira-style starter project, dashboards, API-tester collections, test cases, contract registry |
| [`plugins/ui/`](plugins/ui/) | UI panels: sidebar pages, page slots, Ctrl-K commands, entity tabs — sandboxed iframes |
| [`plugins/links/`](plugins/links/) | Typed cross-entity link types with deep-link URL templates (mock↔GitHub issue, task↔Jira, …) |
| [`plugins/events-and-tasks/`](plugins/events-and-tasks/) | Subscribable notification event types and dispatchable runner task types |
| [`skills/`](skills/) | Agent skills for Claude Code / Cursor / OpenCode: drive Mockarty as an autonomous tester, author plugins end-to-end |

Every plugin folder is self-contained: `plugin.json` (the manifest), assets, and
a README. Folders are also valid `mockarty-cli plugin pack .` inputs — clone,
tweak, repack.

## Migrate from other ecosystems

Already have integrations or mocks elsewhere? `mockarty-cli plugin import`
converts them into installable plugins:

| Source | Command | Becomes |
|---|---|---|
| WireMock mappings | `plugin import ./mappings --from wiremock --id you.stubs` | mock kit |
| Postman collection | `plugin import ./collection.json --from postman --id you.api` | mock kit |
| Mockoon environment | `plugin import ./env.json --from mockoon --id you.env` | mock kit |
| Bruno collection | `plugin import ./my-collection --from bruno --id you.api` | mock kit |
| Insomnia export (v4 JSON/YAML) | `plugin import ./export.json --from insomnia --id you.api` | mock kit |
| Atlassian Connect app | `plugin import ./atlassian-connect.json --from connect` | external UI panels + event types |
| n8n community nodes | `plugin import ./nodes --from n8n` | integration presets |

The [`community.n8n-connectors`](plugins/connectors/n8n-connectors) pack in
this catalogue is the n8n converter's own output over nine real community
nodes — a working reference for what your conversion will look like.

## Publish your own plugin here

1. Author it — the fastest path is `mockarty-cli plugin create my-plugin`
   (see the *Writing plugins* doc; WASM templates include a build script).
2. Validate: `mockarty-cli plugin pack .` then `mockarty-cli plugin inspect <zip>`.
3. Generate the index entry:
   `mockarty-cli plugin publish <zip> --download-url https://github.com/mockarty/extensions/releases/download/<tag>/<zip>`
4. Open a PR adding your plugin folder under the right category and the printed
   object to `index.json`; attach the zip to the release named in the URL.

Note: `index.json` sha256 values must match the exact release assets byte-for-byte — if you rebuild a zip, recompute its sha256 and update the entry.

## Agent skills

The [`skills/`](skills/) folder ships instructions AI assistants load directly:

- **mockarty-drive** — drive Mockarty over MCP as an autonomous backend tester
  (mocks, functional/load/fuzz/chaos/contract runs, TCM, reports).
- **mockarty-plugin-author** — scaffold→build→pack→install a plugin of any
  mechanism, with the WASM guest contract explained.

Install: copy a folder into `.claude/skills/` (Claude Code), `.cursor/rules/`
(Cursor), or reference the `SKILL.md` from any assistant's context. Details in
[`skills/README.md`](skills/README.md).

## License

MIT — see [LICENSE](LICENSE). Plugin folders may carry their own compatible
licenses in their manifests.
