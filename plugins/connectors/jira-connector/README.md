# jira-connector — a Jira plugin, with no code

Shows how to extend Mockarty toward an **external tool** without writing an
adapter. A `connector` contribution is a pre-configured binding onto an
integration adapter Mockarty already ships — here the built-in **Jira** adapter.
Install it and you get a "Jira Cloud" connector with sensible defaults, ready to
turn into a working Jira integration in a couple of clicks.

```json
{
  "contributes": {
    "connectors": [
      { "key": "jira_cloud", "name": "Jira Cloud", "kind": "jira",
        "config_template": { "base_url": "https://your-org.atlassian.net", "project_key": "QA" } }
    ]
  }
}
```

- **kind** — a built-in adapter (`jira`, `github`, `gitlab`, `testrail`, …).
  Validated at install: an unknown kind is refused.
- **config_template** — defaults that pre-fill the integration when created from
  this connector. Never put secrets here — those are entered at setup.

## Install & use

```bash
mockarty-cli plugin pack .
mockarty-cli plugin install mockarty.jira-connector-1.0.0.zip
mockarty-cli plugin enable mockarty.jira-connector
```

Enabled connectors are listed at `GET /api/v1/plugin-connectors` and in the
Plugins UI. Create a Jira integration from the connector — its `kind` and
`config_template` become the defaults; the actual API calls run through
Mockarty's built-in Jira adapter.

This is the pattern for any tracker/CI/chat integration: reuse a ready adapter,
ship a curated connector, distribute it as a plugin.
